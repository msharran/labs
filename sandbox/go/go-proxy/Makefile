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
envoy:
ifeq ($(VALIDATE), 1)
	envoy -c ./envoy-demo.yaml --config-yaml "$$(cat envoy-override.yaml)" --mode validate $(ARGS)
else ifeq ($(DEMO), 1)
	envoy -c ./envoy-demo.yaml --config-yaml "$$(cat envoy-override.yaml)" $(ARGS) 
else
	envoy -c ./envoy.yaml --config-yaml "$$(cat envoy-override.yaml)" $(ARGS) 
endif
