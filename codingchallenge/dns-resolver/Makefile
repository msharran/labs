SHELL := /bin/bash

.PHONY: dns-resolver
dns-resolver: prebuild
	go build -o $@ ./...

.PHONY: prebuild
prebuild:
	go mod tidy
	go fmt ./...
	go vet ./...

