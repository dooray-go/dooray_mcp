BINARY_NAME=dooray
OUTPUT_DIR=dist

.PHONY: build-all clean darwin-amd64 darwin-arm64 linux-amd64 windows-amd64

build-all: darwin-amd64 darwin-arm64 linux-amd64 windows-amd64

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

darwin-amd64: $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME).darwin.amd64 .
	@command -v codesign >/dev/null && codesign --force --sign - $(OUTPUT_DIR)/$(BINARY_NAME).darwin.amd64 || echo "codesign not available, skipping"

darwin-arm64: $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME).darwin.arm64 .
	@command -v codesign >/dev/null && codesign --force --sign - $(OUTPUT_DIR)/$(BINARY_NAME).darwin.arm64 || echo "codesign not available, skipping"

linux-amd64: $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME).linux.amd64 .

windows-amd64: $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME).windows.amd64.exe .

clean:
	rm -rf $(OUTPUT_DIR)
