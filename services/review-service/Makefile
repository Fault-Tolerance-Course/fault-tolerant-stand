
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
	GOBIN=$(LOCAL_BIN) go install github.com/not-for-prod/xo@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.22.1
	GOBIN=$(LOCAL_BIN) go install github.com/ozontech/framer/cmd/framer@v0.1.2


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

fast-generate:
	@$(MAKE) .copy_validate_proto

load-convert:
	$(LOCAL_BIN)/framer convert \
    --from=ozon.json \
    --to=ozon.binary \
    --reflection-addr=localhost:7003 \
    load/requests.json \
    load/requests.json.binary

load-run:
	$(LOCAL_BIN)/framer load  \
	--addr=localhost:7003 \
	--requests-file=load/requests.json.binary \
	--clients 10 \
	const 200 \
	--duration 2m \
	--phout=load/report.phout
