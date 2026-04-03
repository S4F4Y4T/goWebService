BINARY_NAME=server
BUILD_DIR=bin

.PHONY: all build run run-bin clean

all: build

# Run directly from cmd/server.go
run:
	go run cmd/server.go

# Build the binary to bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server.go

# Run the built binary from bin/
run-bin: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
