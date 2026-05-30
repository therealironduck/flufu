[private]
default:
    @just --list

run cmd="":
    #!/usr/bin/env bash
    if [ ! -d "./lib" ]; then
        just prepare
    fi
    YZMA_LIB=./lib go run ./main.go {{ cmd }}

prepare:
    #!/usr/bin/env bash
    go install github.com/hybridgroup/yzma/cmd/yzma@latest
    if [ "$(uname)" = "Darwin" ]; then
        yzma install --lib ./lib --processor metal
    else
        yzma install --lib ./lib --processor cpu
    fi

build:
    go build -o bin/flufu ./main.go

test:
    go test -v ./...
    go vet ./...

lint:
    golangci-lint run ./...
