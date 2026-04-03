BINARY_NAME=server
BUILD_DIR=bin

.PHONY: all build run run-bin dev clean

all: build

dev: run-air

run-air:
	air

# Run directly from cmd/api/maing.go
run:
	go run cmd/api/main.go

# Build the binary to bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/api/main.go

# Run the built binary from bin/
run-bin: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
