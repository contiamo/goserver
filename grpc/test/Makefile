
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


protobuf: $(subst .proto,.pb.go,$(shell find . -name "*.proto")) ## Build the go protobuf files

# get required libraries
.PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
$(.PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

%.pb.go: $(shell find ./ -name "*.proto") $(.PROTOC_GEN_GO)
	@echo "+ protobuf"
	protoc -I. --go_out=plugins=grpc:. *.proto

