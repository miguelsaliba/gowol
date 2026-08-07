.DEFAULT_GOAL := build

.PHONY: fmt vet build build cli server clean run-gowol run-gowol-server

BINDIR := bin

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(BINDIR)/ ./cmd/...

cli: vet
	go build -o $(BINDIR)/gowol ./cmd/gowol

server: vet
	go build -o $(BINDIR)/gowol-server ./cmd/gowol-server

run-gowol:
	go run ./cmd/gowol

run-gowol-server:
	go run ./cmd/gowol-server

release: vet
	CGO_ENABLED=0 go build -ldflags "-s" -o $(BINDIR)/ ./cmd/...

clean:
	rm -rf $(BINDIR)
