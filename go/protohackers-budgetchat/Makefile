SHELL := /bin/bash

build:
	go build -o bin/budgetchat

run:
	go run ./...

generate:
	go generate ./...

TEST_ARGS ?= -v -race
TEST_WHAT ?= ./...
test:
	go test $(TEST_ARGS) $(TEST_WHAT)

installtools:
	go get github.com/stretchr/testify
	go install github.com/vektra/mockery/v2@latest

.PHONY: build test installtools
