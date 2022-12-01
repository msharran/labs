.PHONY: generate, fmt, vet, proto_deps, server, say_hello

serve:
	@go run cmd/server/server.go -port 3456

say_hello:
	@go run cmd/client/client.go -server_addr "localhost:3456"

generate: proto_deps
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/helloworld/helloworld.proto

proto_deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


fmt:
	go fmt ./...

vet:
	go vet ./...
