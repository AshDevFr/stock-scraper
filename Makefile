GOCMD=go
STATIC_PACKCMD=pkger
STATIC_DIR=static
STATIC_MV_CMD=mv web/build $(STATIC_DIR)

YARN_INSTALL=(cd web && yarn)
YARN_BUILD=(cd web && yarn build)

GOBUILD=$(GOCMD) build

BINARY_NAME=stock_scraper
BINARY_DIR=releases
BINARY_DARWIN=$(BINARY_NAME)_darwin
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

.PHONY: all build-deps clean-deps clean build-darwin build-linux build-windows

all: build-darwin build-linux build-windows
	$(MAKE) clean-deps

build-deps:
	$(YARN_INSTALL)
	$(YARN_BUILD)
	$(STATIC_MV_CMD)
	$(STATIC_PACKCMD)

build: clean build-deps
	$(GOBUILD) -o $(BINARY_NAME) -v

clean-deps:
	rm -rf $(STATIC_DIR) || true
	rm -f *-packr.go || true
	rm -rf pkged*.go || true

clean: clean-deps
	rm -r $(BINARY_NAME) || true
	rm -rf $(BINARY_DIR) || true

# Cross compilation
build-darwin: build-deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_DARWIN) -v

build-linux: build-deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_UNIX) -v

build-windows: build-deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)/$(BINARY_WINDOWS) -v
