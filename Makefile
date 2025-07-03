GO_BIN?=$(shell pwd)/.bin
BINARY_DIR=bin
BINARY_NAME=terraform-provider-dotenv

.PHONY: build testacc test lint generate fmt

build:
	mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$(BINARY_NAME)

testacc:
	TF_ACC=1 go test -v ./internal/provider -timeout 5m

test:
	go test -v ./... -timeout 5m

lint:
	golangci-lint run ./...
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx@latest ./...

generate:
	go generate ./...

fmt:
	go fmt ./...
	-go tool gofumpt -l -w .
	-go tool goimports -l -w .
	-go tool wsl -strict-append -test=true -fix ./...
	-go tool perfsprint -fix ./...
	-go tool tfproviderlintx -fix ./...

download:
	@echo Download go.mod dependencies
	@go mod download

tools:
	mkdir -p ${GO_BIN}
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b ${GO_BIN} latest
	GOBIN=${GO_BIN} go install tool
