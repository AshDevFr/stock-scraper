GOCMD=go
STATICPACKCMD=pkger
STATICDIR=static
YARNINSTALL=(cd web && yarn)
YARNBUILD=(cd web && yarn build)
MVSTATICS=mv web/build $(STATICDIR)
GOBUILD=$(GOCMD) build
BINARY_NAME=stock_scraper
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

all: build build-linux build-windows
	$(MAKE) clean-deps

build-deps:
	$(YARNINSTALL)
	$(YARNBUILD)
	$(MVSTATICS)
	$(STATICPACKCMD)

build: clean build-deps
	$(GOBUILD) -o $(BINARY_NAME) -v

clean-deps:
	rm -rf $(STATICDIR)
	rm -f *-packr.go
	rm -rf pkged*.go

clean: clean-deps
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)

# Cross compilation
build-linux: build-deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows: build-deps
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v
