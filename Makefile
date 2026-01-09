.PHONY: build clean test install release all

BINARY_NAME=jwt-decoder
VERSION=1.0.0
BUILD_DIR=build

all: clean build

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) .

install: build
	mv $(BINARY_NAME) /usr/local/bin/

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

# Build for all platforms
release: clean
	mkdir -p $(BUILD_DIR)
	
	# macOS Intel
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	
	# macOS Apple Silicon
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	
	# Linux amd64
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	
	# Linux arm64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	tar -czvf $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-arm64
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	zip -j $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	# Generate checksums
	cd $(BUILD_DIR) && shasum -a 256 *.tar.gz *.zip > checksums.txt

# Run the tool with a sample token
demo:
	go run . eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE3MzYzODAwMDB9.2hDgYvYRtr7VZmHl2XGnM8wLmFaRqW5sJp9aSoYmJBI
