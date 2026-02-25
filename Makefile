.PHONY: build test plugin clean install

# Build the standalone linter
build:
	go build -o bin/loglinter ./cmd/loglinter

# Run tests
test:
	go test -v ./...

# Build as golangci-lint plugin
plugin:
	go build -buildmode=plugin -o loglinter.so ./plugin

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f loglinter.so

# Install the linter
install:
	go install ./cmd/loglinter

# Run golangci-lint with the plugin
lint:
	golangci-lint run --enable=loglinter
