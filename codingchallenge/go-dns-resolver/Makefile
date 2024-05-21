SHELL := /bin/bash

.PHONY: go-dns-resolver
go-dns-resolver: prebuild
	go build -o $@ ./...

.PHONY: prebuild
prebuild:
	go mod tidy
	go fmt ./...
	go vet ./...

