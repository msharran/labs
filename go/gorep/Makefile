# find go files and add them as dependencies
GOFILES := $(shell find . -iname "*.go")

gorep: $(GOFILES)
	go build -o gorep .

ARGS ?= -v -count 1
WHAT ?= ./...
test: fmt vet
	go test $(ARGS) $(WHAT)

fmt: 
	go fmt ./...

vet:
	go vet ./...

print: $(GOFILES)
	@echo $?

install:
	go install

clean: 
	rm -rf gorep

.PHONY: install print clean fmt vet test
