SHELL := /bin/bash

# Find all go files
GOFILES = $(shell find . -iname "*.go")

bin/proxy: $(GOFILES)
	go build -o bin/proxy

.PHONY: proxy
proxy: bin/proxy
	./bin/proxy -downstream_addr localhost:8080


.PHONY: server
server:
	cd ../http; \
	go run .
