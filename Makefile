.PHONY: build
build:
	go build -o bin/mavic ./cmd/mavic/main.go 

.PHONY: build-all
build-all:
	sh ./scripts/build-all.sh ./cmd/mavic/main.go

.PHONY: test
test:
	go test ./...
