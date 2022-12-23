SHELL := /bin/bash

# Find all go files
GOFILES = $(shell find . -iname "*.go")

bin/proxy: $(GOFILES)
	go build -o bin/proxy

.PHONY: proxy
proxy: bin/proxy
	# Add `example.com 127.0.0.1` in /etc/hosts
	./bin/proxy -serverAddr=example.com:8080 -timeout=10 


.PHONY: server
server:
	cd ../http && go run .
