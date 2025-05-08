PKG_VERSION?=0.0.0
COMMIT=`git rev-parse --short HEAD`
PKG_LIST := $(shell go list ./... | grep -v -e /vendor/ -e /api/design -e /scan/logos -e /grpcapi/config/restore_validate )
GO_FILES := $(shell find . -name '*.go' | grep -v -e /vendor/ -e _test.goy)

SERVICE_NAME=load-generation-system

# flags
TEST_FLAGS=-count=1 -race
BUILD_FLAGS=

.PHONY: all init

all: tests build docker

lint: ## Lint the files
	@golangci-lint --timeout=2m run

build: .
	go build ${BUILD_FLAGS} -o ${SERVICE_NAME} cmd/main.go

tests:
	go test ${TEST_FLAGS} ./...

coverage:
	go test ${TEST_FLAGS} -coverprofile=coverage.out ./...
	go tool cover -html="coverage.out"

fmt:
	go fmt ./...

docs-generate:
	go-swagger3 --module-path . --main-file-path cmd/main.go --handler-path api/rest/manager/handlers --output .swag/swagger.json --schema-without-pkg

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./api/inject/

run:
	go run cmd/main.go manager

run-node:
	go run cmd/main.go node