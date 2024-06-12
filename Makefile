BIN_NAME=todayornever

.PHONY: build

all: build

build:
	@mkdir -p build
	go build -o build/$(BIN_NAME)

run:
	./build/$(BIN_NAME)

clean:
	rm -rf ./build ./tmp

build-all:
	@mkdir -p build
	GOOS=windows GOARCH=amd64 go build -o build/$(BIN_NAME)-windows-amd64.exe # Windows 64bits
	GOOS=darwin GOARCH=amd64 go build -o build/$(BIN_NAME)-darwin-amd64 # MacOS Intel
	GOOS=darwin GOARCH=arm64 go build -o build/$(BIN_NAME)-darwin-arm64 # MacOS Silicon
	GOOS=linux GOARCH=amd64 go build -o build/$(BIN_NAME)-linux-amd64 # Linux 64bits
