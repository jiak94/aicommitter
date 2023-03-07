# Define variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DEST=bin
CREATE_DIR=mkdir $(DEST)

# Define flags
BUILD_FLAGS=-v -ldflags="-s -w"

# Define targets
.PHONY: all clean test

all: clean test build

build: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(DEST)/aicommit-linux-amd64

build-linux-arm64:
	env GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(DEST)/aicommit-linux-arm64

build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(DEST)/aicommit-darwin-amd64

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(DEST)/aicommit-darwin-arm64

clean:
	$(GOCLEAN)
	rm -rf $(DEST)
