// Fetch wallpaper-quality images from NASA mission galleries.
//
// Sources:
//   - science.nasa.gov content-list API (Cassini)
//   - Flickr album pages              (JWST, Hubble)
//
// Usage:
//
//	go run getPhotos.go [--limit N] [--out-dir DIR] [--debug]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ---------------------------------------------------------------------------
// Resolution / ratio constants
// ---------------------------------------------------------------------------

const (
	minWidth  = 1920
	minHeight = 1080
	minRatio  = 1.4
	maxRatio  = 2.4
)

// ---------------------------------------------------------------------------
// Concurrency knobs
// ---------------------------------------------------------------------------

const (
	apiRate    = 2.5 // max API/scrape requests per second (shared globally)
	apiWorkers = 6   // max concurrent in-flight API requests
	dlWorkers  = 8   // max concurrent image downloads
)

// ---------------------------------------------------------------------------
// Mission definitions
// ---------------------------------------------------------------------------

type flickrAlbum struct {
	albumID   string
	userPath  string // flickr username in URL, e.g. "nasahubble"
	subfolder string // subdirectory inside mission dir; empty = flat
}

type missionDef struct {
	name        string
	scienceSlug string        // non-empty → use NASA Science scraper
	flickr      []flickrAlbum // non-empty → use Flickr scraper
}

var missions = []missionDef{
	{
		name:        "cassini",
		scienceSlug: "cassini",
	},
	{
		name: "jwst",
		flickr: []flickrAlbum{
			{albumID: "72177720332131144", userPath: "nasawebbtelescope"},
			{albumID: "72177720332633211", userPath: "nasawebbtelescope"},
		},
	},
	{
		name: "hubble",
		flickr: []flickrAlbum{
			{albumID: "72157661867344528", userPath: "nasahubble", subfolder: "galaxies"},
			{albumID: "72157695205167691", userPath: "nasahubble", subfolder: "nebula"},
			{albumID: "72157698290709165", userPath: "nasahubble", subfolder: "star-clusters"},
		},
	},
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

type imageRecord struct {
	Title    string
	ImageURL string
	Width    int
	Height   int
	PiaID    string
}

// ---------------------------------------------------------------------------
// Rate limiter — token bucket, one instance shared across all goroutines
// ---------------------------------------------------------------------------

type rateLimiter struct{ tokens chan struct{} }

func newRateLimiter(rps float64, burst int) *rateLimiter {
	r := &rateLimiter{tokens: make(chan struct{}, burst)}
	for i := 0; i < burst; i++ {
		r.tokens <- struct{}{}
	}
	go func() {
		t := time.NewTicker(time.Duration(float64(time.Second) / rps))
		defer t.Stop()
		for range t.C {
			select {
			case r.tokens <- struct{}{}:
			default:
			}
		}
	}()
	return r
}

func (r *rateLimiter) wait() { <-r.tokens }

// ---------------------------------------------------------------------------
// Globals
// ---------------------------------------------------------------------------

var (
	apiLimiter *rateLimiter
	client     = &http.Client{Timeout: 30 * time.Second}
	debugMode  bool

	imagePostTypes = map[string]bool{
		"image-article-ext": true,
		"image-article":     true,
	}
	piaRe = regexp.MustCompile(`(?i)\bPIA\d{4,6}\b`)
	tagRe = regexp.MustCompile(`<[^>]+>`)
)

// ---------------------------------------------------------------------------
// HTTP / text helpers
// ---------------------------------------------------------------------------

func httpGet(rawURL string) ([]byte, error) {
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
	}
	return io.ReadAll(resp.Body)
}

func stripTags(s string) string {
	return strings.TrimSpace(tagRe.ReplaceAllString(s, ""))
}

func findPIA(texts ...string) string {
	for _, t := range texts {
		if m := piaRe.FindString(t); m != "" {
			return strings.ToUpper(m)
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Resolution filter
// ---------------------------------------------------------------------------

func passesFilter(w, h int) bool {
	if w < minWidth || h < minHeight {
		return false
	}
	ratio := float64(w) / float64(h)
	return ratio >= minRatio && ratio <= maxRatio
}

// ---------------------------------------------------------------------------
// NASA Science scraper
// ---------------------------------------------------------------------------

// galleryConfig mirrors the data-starting-query JSON block.
// requesting_id can be a JSON number or string; we store it as RawMessage
// and stringify it when needed. base_terms is a JSON-encoded string value.
type galleryConfig struct {
	RequestingID json.RawMessage `json:"requesting_id"`
	BaseTerms    string          `json:"base_terms"`
	BlockID      string          `json:"block_id"`
	PostTypes    []string        `json:"post_types"`
}

func (c *galleryConfig) requestingIDStr() string {
	// Strip surrounding quotes if it was encoded as a JSON string
	return strings.Trim(string(c.RequestingID), `"`)
}

func discoverGalleryConfig(missionSlug string) (*galleryConfig, error) {
	pageURL := fmt.Sprintf("https://science.nasa.gov/mission/%s/multimedia/images/", missionSlug)
	if debugMode {
		fmt.Printf("  [debug] gallery page: %s\n", pageURL)
	}
	apiLimiter.wait()
	body, err := httpGet(pageURL)
	if err != nil {
		return nil, fmt.Errorf("fetch gallery page: %w", err)
	}
	re := regexp.MustCompile(`data-starting-query="([^"]+)"`)
	for _, m := range re.FindAllSubmatch(body, -1) {
		raw := html.UnescapeString(string(m[1]))
		var cfg galleryConfig
		if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
			continue
		}
		if len(cfg.RequestingID) > 0 && cfg.BaseTerms != "" {
			if debugMode {
				fmt.Printf("  [debug] config: requesting_id=%s\n", cfg.requestingIDStr())
			}
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("gallery config not found for '%s'", missionSlug)
}

type contentItem struct {
	PostType  string `json:"postType"`
	Permalink string `json:"permalink"`
}

func fetchContentPage(cfg *galleryConfig, page, perPage int) ([]contentItem, error) {
	postTypes := cfg.PostTypes
	if len(postTypes) == 0 {
		postTypes = []string{"resource"}
	}
	ptJSON, _ := json.Marshal(postTypes)

	params := url.Values{
		"block_id":        {cfg.BlockID},
		"requesting_id":   {cfg.requestingIDStr()},
		"post_types":      {string(ptJSON)},
		"base_terms":      {cfg.BaseTerms},
		"number_of_items": {fmt.Sprintf("%d", perPage)},
		"current_page":    {fmt.Sprintf("%d", page)},
		"response_format": {"json"},
	}
	reqURL := "https://science.nasa.gov/wp-json/smd/v1/content-list/?" + params.Encode()
	if debugMode {
		fmt.Printf("  [debug] content-list page %d\n", page)
	}
	apiLimiter.wait()
	body, err := httpGet(reqURL)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Type struct {
			Content []contentItem `json:"content"`
		} `json:"type"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Type.Content, nil
}

func fetchNASAImageDetail(permalink string) (*imageRecord, error) {
	slug := strings.TrimSuffix(permalink, "/")
	if i := strings.LastIndex(slug, "/"); i >= 0 {
		slug = slug[i+1:]
	}
	apiURL := fmt.Sprintf("https://www.nasa.gov/wp-json/wp/v2/image-article?slug=%s&_embed",
		url.QueryEscape(slug))
	if debugMode {
		fmt.Printf("    [debug] detail: %s\n", apiURL)
	}
	apiLimiter.wait()
	body, err := httpGet(apiURL)
	if err != nil {
		return nil, err
	}
	var items []struct {
		Title struct{ Rendered string `json:"rendered"` } `json:"title"`
		Embedded struct {
			FeaturedMedia []struct {
				SourceURL    string `json:"source_url"`
				Caption      struct{ Rendered string `json:"rendered"` } `json:"caption"`
				MediaDetails struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"media_details"`
			} `json:"wp:featuredmedia"`
		} `json:"_embedded"`
	}
	if err := json.Unmarshal(body, &items); err != nil || len(items) == 0 {
		return nil, err
	}
	item := items[0]
	if len(item.Embedded.FeaturedMedia) == 0 {
		return nil, nil
	}
	media := item.Embedded.FeaturedMedia[0]
	if media.SourceURL == "" {
		return nil, nil
	}
	title := stripTags(item.Title.Rendered)
	caption := stripTags(media.Caption.Rendered)
	return &imageRecord{
		Title:    title,
		ImageURL: media.SourceURL,
		Width:    media.MediaDetails.Width,
		Height:   media.MediaDetails.Height,
		PiaID:    findPIA(media.SourceURL, title, caption),
	}, nil
}

func scrapeNASAScience(missionSlug string, limit int) ([]imageRecord, error) {
	cfg, err := discoverGalleryConfig(missionSlug)
	if err != nil {
		return nil, err
	}

	// Collect image-article permalinks from paginated content-list
	var permalinks []string
	seen := map[string]bool{}
	for page := 1; ; page++ {
		items, err := fetchContentPage(cfg, page, 100)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  content-list page %d failed: %v\n", page, err)
			break
		}
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			if imagePostTypes[item.PostType] && item.Permalink != "" && !seen[item.Permalink] {
				seen[item.Permalink] = true
				permalinks = append(permalinks, item.Permalink)
			}
		}
		if len(items) < 100 {
			break
		}
	}
	if limit > 0 && len(permalinks) > limit {
		permalinks = permalinks[:limit]
	}

	// Fetch all detail pages concurrently
	type result struct {
		rec *imageRecord
		err error
	}
	results := make([]result, len(permalinks))
	var wg sync.WaitGroup
	sem := make(chan struct{}, apiWorkers)
	for i, pl := range permalinks {
		wg.Add(1)
		go func(i int, pl string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			rec, err := fetchNASAImageDetail(pl)
			results[i] = result{rec, err}
		}(i, pl)
	}
	wg.Wait()

	var records []imageRecord
	for _, r := range results {
		if r.err != nil {
			if debugMode {
				fmt.Fprintf(os.Stderr, "  detail error: %v\n", r.err)
			}
			continue
		}
		if r.rec != nil {
			records = append(records, *r.rec)
		}
	}
	return records, nil
}

// ---------------------------------------------------------------------------
// Flickr scraper
// ---------------------------------------------------------------------------

var flickrPhotoIDRe = regexp.MustCompile(`"id":"(\d{10,})"`)

// scrapeFlickrAlbum fetches the album page and extracts photo IDs,
// then fetches each photo's page to get the best available image URL.
func scrapeFlickrAlbum(albumID, userPath string, limit int) ([]imageRecord, error) {
	albumURL := fmt.Sprintf("https://www.flickr.com/photos/%s/albums/%s/", userPath, albumID)
	if debugMode {
		fmt.Printf("  [debug] flickr album page: %s\n", albumURL)
	}
	apiLimiter.wait()
	body, err := httpGet(albumURL)
	if err != nil {
		return nil, fmt.Errorf("fetch album page: %w", err)
	}

	// Extract unique photo IDs (album page embeds them server-side)
	var photoIDs []string
	seen := map[string]bool{albumID: true}
	for _, m := range flickrPhotoIDRe.FindAllSubmatch(body, -1) {
		id := string(m[1])
		if !seen[id] {
			seen[id] = true
			photoIDs = append(photoIDs, id)
		}
	}
	if debugMode {
		fmt.Printf("  [debug] flickr: %d photo IDs from album page\n", len(photoIDs))
	}
	if limit > 0 && len(photoIDs) > limit {
		photoIDs = photoIDs[:limit]
	}

	// Fetch each photo page concurrently to find best size URL
	type result struct {
		rec *imageRecord
		err error
	}
	results := make([]result, len(photoIDs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, apiWorkers)
	for i, id := range photoIDs {
		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			rec, err := fetchFlickrPhotoDetail(id, userPath)
			results[i] = result{rec, err}
		}(i, id)
	}
	wg.Wait()

	var records []imageRecord
	for _, r := range results {
		if r.err != nil {
			if debugMode {
				fmt.Fprintf(os.Stderr, "  flickr detail error: %v\n", r.err)
			}
			continue
		}
		if r.rec != nil {
			records = append(records, *r.rec)
		}
	}
	return records, nil
}

type flickrSize struct {
	DisplayURL string `json:"displayUrl"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

// fetchFlickrPhotoDetail fetches the Flickr photo page HTML and extracts
// the largest available image URL from the embedded sizes JSON object.
func fetchFlickrPhotoDetail(photoID, userPath string) (*imageRecord, error) {
	photoURL := fmt.Sprintf("https://www.flickr.com/photos/%s/%s/", userPath, photoID)
	if debugMode {
		fmt.Printf("    [debug] flickr photo: %s\n", photoURL)
	}
	apiLimiter.wait()
	body, err := httpGet(photoURL)
	if err != nil {
		return nil, err
	}

	// The photo page embeds: "sizes":{"sq":{...},"q":{...},...}
	// Find the opening brace for the sizes object and extract it
	idx := strings.Index(string(body), `"sizes":{`)
	if idx < 0 {
		return nil, fmt.Errorf("no sizes data in page for photo %s", photoID)
	}
	chunk := string(body)[idx+len(`"sizes":`):]
	depth, end := 0, 0
	for i, c := range chunk {
		if c == '{' {
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}
	if end == 0 {
		return nil, fmt.Errorf("could not parse sizes JSON for photo %s", photoID)
	}

	var sizes map[string]flickrSize
	if err := json.Unmarshal([]byte(chunk[:end]), &sizes); err != nil {
		return nil, fmt.Errorf("sizes JSON parse error for %s: %w", photoID, err)
	}

	// Pick the size with the greatest width
	var best flickrSize
	for _, s := range sizes {
		if s.Width > best.Width {
			best = s
		}
	}
	if best.DisplayURL == "" {
		return nil, fmt.Errorf("no sizes found for photo %s", photoID)
	}

	// DisplayURL uses protocol-relative "//..." form
	imgURL := best.DisplayURL
	if strings.HasPrefix(imgURL, "//") {
		imgURL = "https:" + imgURL
	}

	return &imageRecord{
		Title:    photoID,
		ImageURL: imgURL,
		Width:    best.Width,
		Height:   best.Height,
	}, nil
}

// ---------------------------------------------------------------------------
// Download helpers
// ---------------------------------------------------------------------------

func downloadFile(rawURL, dest string) error {
	req, _ := http.NewRequest("GET", rawURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// checkResolution uses image.DecodeConfig to read only the image header —
// much faster than decoding the full pixel data.
func checkResolution(path string) (ok bool, w, h int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return false, 0, 0, err
	}
	return passesFilter(cfg.Width, cfg.Height), cfg.Width, cfg.Height, nil
}

// ---------------------------------------------------------------------------
// Download records into a target directory (concurrent, bounded by dlWorkers)
// ---------------------------------------------------------------------------

func downloadRecords(label string, records []imageRecord, dir string) int {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "  [%s] mkdir %s failed: %v\n", label, dir, err)
		return 0
	}

	existing, _ := os.ReadDir(dir)
	existingNames := make(map[string]bool, len(existing))
	for _, e := range existing {
		existingNames[e.Name()] = true
	}

	type dlTask struct {
		rec    imageRecord
		safeID string
		ext    string
		dest   string
	}

	var tasks []dlTask
	count := 0

	for _, rec := range records {
		rawURL := strings.SplitN(rec.ImageURL, "?", 2)[0]
		base := rawURL[strings.LastIndex(rawURL, "/")+1:]
		extIdx := strings.LastIndex(base, ".")
		var safeID, ext string
		if extIdx >= 0 {
			safeID, ext = base[:extIdx], base[extIdx:]
		} else {
			safeID, ext = base, ".jpg"
		}
		if rec.PiaID != "" {
			safeID = rec.PiaID
		}
		safeID = strings.NewReplacer("/", "_", " ", "_").Replace(safeID)
		if safeID == "" || rec.ImageURL == "" {
			continue
		}

		// Skip if any file with this base name already exists
		alreadyExists := false
		for fn := range existingNames {
			if strings.HasPrefix(fn, safeID) {
				alreadyExists = true
				break
			}
		}
		if alreadyExists {
			fmt.Printf("  [%s] Already exists: %s\n", label, safeID)
			count++
			continue
		}

		// Pre-filter using metadata dimensions (avoids downloading obvious rejects)
		if rec.Width > 0 && rec.Height > 0 && !passesFilter(rec.Width, rec.Height) {
			if debugMode {
				fmt.Printf("  [debug] [%s] Skip %s (%dx%d fails filter)\n", label, safeID, rec.Width, rec.Height)
			}
			continue
		}

		tasks = append(tasks, dlTask{
			rec:    rec,
			safeID: safeID,
			ext:    ext,
			dest:   filepath.Join(dir, safeID+ext),
		})
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, dlWorkers)

	for _, t := range tasks {
		wg.Add(1)
		go func(t dlTask) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			wStr, hStr := "?", "?"
			if t.rec.Width > 0 {
				wStr = fmt.Sprintf("%d", t.rec.Width)
			}
			if t.rec.Height > 0 {
				hStr = fmt.Sprintf("%d", t.rec.Height)
			}
			fmt.Printf("  [%s] Downloading %s%s  (%sx%s)\n", label, t.safeID, t.ext, wStr, hStr)

			if err := downloadFile(t.rec.ImageURL, t.dest); err != nil {
				fmt.Fprintf(os.Stderr, "  [%s] Download failed %s: %v\n", label, t.safeID, err)
				return
			}

			ok, w, h, err := checkResolution(t.dest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [%s] Warning: cannot check dimensions of %s: %v\n", label, t.safeID, err)
				mu.Lock()
				count++
				mu.Unlock()
				return
			}
			if !ok {
				os.Remove(t.dest)
				ratio := float64(w) / float64(h)
				if w < minWidth || h < minHeight {
					fmt.Printf("  [%s] Removed %s (too small: %dx%d)\n", label, t.safeID, w, h)
				} else {
					fmt.Printf("  [%s] Removed %s (bad ratio: %.2f, want %.1f–%.1f)\n",
						label, t.safeID, ratio, minRatio, maxRatio)
				}
				return
			}
			mu.Lock()
			count++
			mu.Unlock()
		}(t)
	}
	wg.Wait()
	return count
}

// ---------------------------------------------------------------------------
// Mission runners
// ---------------------------------------------------------------------------

func runMission(def missionDef, outDir string, limit int) int {
	missionDir := filepath.Join(outDir, def.name)

	if def.scienceSlug != "" {
		// NASA Science scraper
		fmt.Printf("[%s] Scraping science.nasa.gov/mission/%s/...\n", def.name, def.scienceSlug)
		records, err := scrapeNASAScience(def.scienceSlug, limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[%s] scrape failed: %v\n", def.name, err)
			return 0
		}
		fmt.Printf("[%s] %d records retrieved\n", def.name, len(records))
		n := downloadRecords(def.name, records, missionDir)
		fmt.Printf("[%s] %d images saved to %s\n", def.name, n, missionDir)
		return n
	}

	// Flickr scraper — one goroutine per album, all run concurrently
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		total int
	)
	for _, album := range def.flickr {
		wg.Add(1)
		go func(album flickrAlbum) {
			defer wg.Done()
			label := def.name
			dir := missionDir
			if album.subfolder != "" {
				label = fmt.Sprintf("%s/%s", def.name, album.subfolder)
				dir = filepath.Join(missionDir, album.subfolder)
			}
			fmt.Printf("[%s] Scraping Flickr album %s...\n", label, album.albumID)
			records, err := scrapeFlickrAlbum(album.albumID, album.userPath, limit)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] album %s scrape failed: %v\n", label, album.albumID, err)
				return
			}
			fmt.Printf("[%s] %d records from album %s\n", label, len(records), album.albumID)
			n := downloadRecords(label, records, dir)
			fmt.Printf("[%s] %d images saved to %s\n", label, n, dir)
			mu.Lock()
			total += n
			mu.Unlock()
		}(album)
	}
	wg.Wait()
	return total
}

// ---------------------------------------------------------------------------
// Entry point — all missions run concurrently
// ---------------------------------------------------------------------------

func main() {
	var limit int
	var outDir string
	flag.IntVar(&limit, "limit", 10, "Images to fetch per mission/album (default: 10)")
	flag.StringVar(&outDir, "out-dir", ".", "Root directory for saved images")
	flag.BoolVar(&debugMode, "debug", false, "Print scraper diagnostics")
	flag.Parse()

	apiLimiter = newRateLimiter(apiRate, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 0

	for _, def := range missions {
		wg.Add(1)
		go func(def missionDef) {
			defer wg.Done()
			n := runMission(def, outDir, limit)
			mu.Lock()
			total += n
			mu.Unlock()
		}(def)
	}
	wg.Wait()

	fmt.Printf("\nDone. %d total images across %d missions.\n", total, len(missions))
}
