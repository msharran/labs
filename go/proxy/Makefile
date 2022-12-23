SHELL := /bin/bash

# Find all go files
GOFILES = $(shell find . -iname "*.go")

bin/proxy: $(GOFILES)
	go build -o bin/proxy

.PHONY: proxy
proxy: bin/proxy
	./bin/proxy -serverAddr=localhost:8080 -timeout=10 


.PHONY: httpserver
httpserver:
	cd httpserver && go run .


.PHONY: envoy
ENVOY := envoy -c ./envoy-demo.yaml --config-yaml "$$(cat envoy-override.yaml)" 
envoy:
ifeq ($(VALIDATE), 1)
	$(ENVOY) --mode validate $(ARGS) 
else
	$(ENVOY) $(ARGS) 
endif
