.DEFAULT_GOAL := run
.PHONY: fmt vet build run test

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@go build -o ./bin/cli ./cmd/cli/main.go

run: build 
	@go run ./cmd/cli/main.go

test:
	@go test ./...



