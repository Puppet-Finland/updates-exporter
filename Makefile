# Makefile

# 1. Grab the latest git tag. If no tags exist, fall back to the short commit hash.
VERSION ?= $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "unknown")

# 2. Define the full path to the variable we want to inject into.
# Format: package_path.variable_name
VERSION_FLAG_PATH = main.Version

# 3. Combine into ldflags
LDFLAGS = -ldflags "-X $(VERSION_FLAG_PATH)=$(VERSION)"

.PHONY: build clean

# The default target
build:
	@echo "Building updates-exporter with version: $(VERSION)"
	go build $(LDFLAGS) -o bin/updates-exporter main.go

clean:
	rm -f bin/*
