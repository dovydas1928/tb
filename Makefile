# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_SERVER_NAME=tb-server
BINARY_CLIENT_NAME=tb-client
PROTO_GEN_GO=$(shell which protoc-gen-go)
PROTO_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)

# Check for protoc-gen-go and protoc-gen-go-grpc plugins installation
ifeq ($(PROTO_GEN_GO),)
	$(error "protoc-gen-go is not installed. Please run 'go install google.golang.org/protobuf/cmd/protoc-gen-go@latest'")
endif

ifeq ($(PROTO_GEN_GO_GRPC),)
	$(error "protoc-gen-go-grpc is not installed. Please run 'go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'")
endif

# gRPC
PROTO_FILES=./pkg/proto/*.proto
PROTO_OUT_DIR=./pkg/proto

all: build
build: server client
proto:
	@echo "Generating Go files from proto..."
	protoc --proto_path=$(PROTO_OUT_DIR) --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
			--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
			$(PROTO_FILES)

server: proto
	@echo "Building server..."
	$(GOBUILD) -o $(BINARY_SERVER_NAME) ./cmd/tb-server

client: proto
	@echo "Building client..."
	$(GOBUILD) -o $(BINARY_CLIENT_NAME) ./cmd/tb-client

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_SERVER_NAME)
	rm -f $(BINARY_CLIENT_NAME)
	rm -f $(PROTO_OUT_DIR)/*.go

.PHONY: all build server client clean proto
