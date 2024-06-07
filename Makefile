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
