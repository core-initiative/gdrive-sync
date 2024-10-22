.PHONY: all build run run-immediate test clean

# Define variables
BINARY_NAME=drive-sync
CONFIG_FILE=config.yaml

# Default target
all: build

# Build the binary
build:
	go build -o build/$(BINARY_NAME)

# Run the application with the scheduled configuration
run: build
	./build/$(BINARY_NAME) -config ./conf/$(CONFIG_FILE)

# Run the application immediately for testing
run-immediate: build
	./build/$(BINARY_NAME) -config ./conf/$(CONFIG_FILE) -immediate

# Test the application
test:
	go test ./...

# Clean the build
clean:
	go clean
	rm -f ./build/$(BINARY_NAME)