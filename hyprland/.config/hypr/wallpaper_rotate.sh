#!/usr/bin/env bash
# Rotate wallpapers from ~/.config/background_photos/ every 5 minutes

PHOTOS_DIR="$HOME/.config/background_photos"
INTERVAL=300  # 5 minutes in seconds
MONITORS=("DP-3" "HDMI-A-1")
LOG="$HOME/.local/share/hypr/wallpaper_rotate.log"

mkdir -p "$(dirname "$LOG")"
exec >> "$LOG" 2>&1

echo "[$(date)] wallpaper_rotate.sh started"

# Wait for swww-daemon to be ready
sleep 2

get_random_image() {
    find -L "$PHOTOS_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.tif" -o -iname "*.tiff" \) | shuf -n 1
}

while true; do
    for monitor in "${MONITORS[@]}"; do
        img=$(get_random_image)
        if [[ -n "$img" ]]; then
            echo "[$(date)] Setting $monitor → $img"
            awww img --outputs "$monitor" --transition-type fade "$img" || echo "[$(date)] ERROR: awww failed for $monitor"
        else
            echo "[$(date)] WARNING: No images found in $PHOTOS_DIR"
        fi
    done
    sleep "$INTERVAL"
done
