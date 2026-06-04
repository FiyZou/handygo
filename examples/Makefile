.PHONY: help install-tools generate gen-model go-generate run dev test build clean smoke

HANDYGO ?= handygo
APP_CONFIG ?= manifest/config.local.yaml
BIN ?= handygo-example

help:
	@printf '%s\n' 'Targets:'
	@printf '%s\n' '  install-tools  Install the handygo CLI'
	@printf '%s\n' '  generate       Generate models from manifest/gen.yaml'
	@printf '%s\n' '  run            Run with embedded manifest/config.yaml'
	@printf '%s\n' '  dev            Run with APP_CONFIG (default manifest/config.local.yaml)'
	@printf '%s\n' '  test           Run tests'
	@printf '%s\n' '  build          Build the example binary'
	@printf '%s\n' '  smoke          Check the local health endpoint'
	@printf '%s\n' '  clean          Remove local build/database artifacts'

install-tools:
	go install github.com/FiyZou/handygo/cmd/handygo@latest

generate: gen-model

gen-model:
	$(HANDYGO) gen model -c manifest/gen.yaml

go-generate:
	go generate ./...

run:
	go run .

dev:
	APP_CONFIG=$(APP_CONFIG) go run .

test:
	go test ./...

build:
	go build -o $(BIN) .

smoke:
	curl -fsS http://127.0.0.1:8080/health

clean:
	rm -f $(BIN) handygo-example.db handygo-example.db-shm handygo-example.db-wal
