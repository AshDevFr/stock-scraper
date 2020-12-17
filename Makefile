# Go parameters
GOCMD=go
STATICPACKCMD=pkger
YARNBUILD=(cd web && yarn build)
GOBUILD=$(GOCMD) build
BINARY_NAME=stock_scraper
BINARY_UNIX=$(BINARY_NAME)
BINARY_WINDOWS=$(BINARY_NAME).exe

all: build
build-deps:
	$(YARNBUILD)
	$(STATICPACKCMD)

build:
	$(YARNBUILD)
	$(STATICPACKCMD)
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f *-packr.go
	rm -rf pkged*.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v
