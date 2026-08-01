.DEFAULT_GOAL := build

.PHONY: fmt vet build build cli server clean run-gowol run-gowol-server

BINDIR := bin

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet cli server

cli: vet
	go build -o $(BINDIR)/gowol ./cmd/gowol

server: vet
	go build -o $(BINDIR)/gowol-server ./cmd/gowol-server

run-gowol:
	go run ./cmd/gowol

run-gowol-server:
	go run ./cmd/gowol-server

clean:
	rm -rf $(BINDIR)
