.PHONY: all generate generate-proto generate-swagger build test run-screen run-content clean

# Binary paths
SWAG := $(shell go env GOPATH)/bin/swag

export DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/signage?sslmode=disable
export ENV ?= development

all: generate build test

generate: generate-proto generate-swagger

generate-proto:
	@echo "Generating gRPC proto files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/v1/screen/*.proto api/proto/v1/content/*.proto

generate-swagger:
	@echo "Generating Swagger documents..."
	$(SWAG) init -g cmd/screen-service/main.go

build:
	@echo "Building services..."
	go build -o bin/screen-service cmd/screen-service/main.go
	go build -o bin/content-service cmd/content-service/main.go

test:
	@echo "Running tests..."
	go test -v ./...

run-screen:
	@echo "Running screen service..."
	go run cmd/screen-service/main.go

run-content:
	@echo "Running content service..."
	go run cmd/content-service/main.go

clean:
	@echo "Cleaning up..."
	rm -rf bin/
