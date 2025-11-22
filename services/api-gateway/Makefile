
# Директория, в которой хранятся исполняемые
# файлы проекта и зависимости, необходимые для сборки.
LOCAL_BIN := $(CURDIR)/bin

export PATH := $(PATH):$(LOCAL_BIN)

# ==================================================================================== #
# PREPARE
# ==================================================================================== #

bin-deps:
	@$(MAKE) dependency

.PHONY: dependency
dependency:
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@v1.6.0
	GOBIN=$(LOCAL_BIN) go install gotest.tools/gotestsum@v1.7.0

.PHONY: run
run:
	go run cmd/main.go

.PHONY: .go-gen
.go-gen:
	go generate ./internal/...

generate:
	@$(MAKE) .go-gen
	buf dep update
	buf dep prune
	buf build
	buf generate --template buf.gen.local.yaml
	buf generate --template buf.gen.external.yaml


fast-generate:
	@$(MAKE) .copy_validate_proto
