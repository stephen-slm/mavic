.PHONY: build
build:
	go build -o bin/mavic ./cmd/mavic/main.go 

.PHONY: test
test:
	go test ./...
