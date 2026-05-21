# Makefile

# 1. Grab the latest git tag. If no tags exist, fall back to the short commit hash.
APP_VERSION ?= $(shell git describe --tags --always --abbrev=0 2>/dev/null || echo "unknown")

# 3. Combine into ldflags
LDFLAGS = -ldflags "-X main.Version=$(APP_VERSION)"

.PHONY: build clean

# The default target
build:
	@mkdir -p bin
	@echo "Building updates-exporter with version: $(APP_VERSION)"
	go build $(LDFLAGS) -o bin/updates-exporter main.go

clean:
	rm -f bin/*
