.PHONY: build test clean install

# Build the binary
build:
	go build

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f cpcommit

# Install to GOPATH/bin
install:
	go install

# Run go mod tidy
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...
