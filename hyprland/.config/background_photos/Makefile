BINARY   := getPhotos
OUT_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
LIMIT    := 10
MISSIONS := cassini hubble jwst new-horizons

.PHONY: all run clean-photos clean-binary clean help

all: $(BINARY)

$(BINARY): getPhotos.go go.mod go.sum
	go build -o $(BINARY) .

## run: build (if needed) and fetch photos
run: $(BINARY)
	./$(BINARY) --limit $(LIMIT) --out-dir $(OUT_DIR)

## clean-photos: delete all downloaded mission photo directories
clean-photos:
	@echo "Removing downloaded photos..."
	rm -rf $(addprefix $(OUT_DIR), $(MISSIONS))

## clean-binary: remove the compiled binary
clean-binary:
	rm -f $(BINARY)

## clean: remove photos and binary
clean: clean-photos clean-binary

## help: list available targets
help:
	@grep -E '^##' Makefile | sed 's/^## /  /'
