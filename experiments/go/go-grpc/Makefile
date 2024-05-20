.PHONY: generate, fmt, vet, proto_deps, server, say_hello

PORT ?= 3456
ADDR ?= "localhost:$(PORT)"

build:
	go build -o bin/server cmd/server/server.go
	go build -o bin/client cmd/client/client.go

DRY_RUN ?= 0
ifeq ($(DRY_RUN), 1)
DRY_RUN_FLAG = --dry-run=server
else
DRY_RUN_FLAG =
endif

# find directories with kustomization.yaml files
KUSTOMIZE_TARGETS = $(shell find kustomize/overlays -name kustomization.yaml -exec dirname {} \;)
ACTION ?= build

.PHONY: $(KUSTOMIZE_TARGETS)
$(KUSTOMIZE_TARGETS):
ifeq ($(ACTION), build)
	kustomize build $@
else ifeq ($(ACTION), apply)
	kustomize build $@ | kubectl apply --server-side $(DRY_RUN_FLAG) -f -
else ifeq ($(ACTION), destroy)
	kustomize build $@ | kubectl delete $(DRY_RUN_FLAG) -f -
endif

.PHONY: helm-install
helm-install:
	helm install --generate-name helm

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(shell helm list -q)

.PHONY: server
serve:
	@go run cmd/server/server.go -port $(PORT)

.PHONY: generate
generate: proto_deps
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/helloworld/helloworld.proto

.PHONY: proto_deps
proto_deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...
