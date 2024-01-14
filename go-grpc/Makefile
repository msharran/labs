.PHONY: generate, fmt, vet, proto_deps, server, say_hello

PORT ?= 3456
ADDR ?= "localhost:$(PORT)"

build:
	go build -o bin/server cmd/server/server.go
	go build -o bin/client cmd/client/client.go

serve:
	@go run cmd/server/server.go -port $(PORT)

say_hello:
	@go run cmd/client/client.go -server_addr $(ADDR)

generate: proto_deps
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/helloworld/helloworld.proto

proto_deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


fmt:
	go fmt ./...

vet:
	go vet ./...
