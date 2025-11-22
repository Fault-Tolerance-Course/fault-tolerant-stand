
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

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

DB_NAME ?= payments
DB_USER ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
MIGRATION_FOLDER=./tools/migrations

.PHONY: migrations-up ## накатывает миграции на базу данных
migrations-up:
	$(LOCAL_BIN)/goose postgres 'host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) sslmode=disable dbname=$(DB_NAME)' -dir $(MIGRATION_FOLDER) -allow-missing up

.PHONY: migrations-reset ## накатывает миграции на базу данных
migrations-reset:
	$(LOCAL_BIN)/goose postgres 'host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) sslmode=disable dbname=$(DB_NAME)' -dir $(MIGRATION_FOLDER) -allow-missing reset

.PHONY: migrations ## накатывает миграции на базу данных
migrations: migrations-reset migrations-up


# ==================================================================================== #
# CODE GENERATION
# ==================================================================================== #

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


XO_OUTPUT_PATH=./internal/pkg/xo
XO_TEMPLATE_PATH=./tools/xo_templates
.PHONY: xo ## генерация dto базы данных
xo: migrations-up
	mkdir -p $(XO_OUTPUT_PATH)
	xo "pgsql://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
	-o $(XO_OUTPUT_PATH) --template-path $(XO_TEMPLATE_PATH) --schema public --suffix ".xo.go" --custom-type-package custom
	rm $(XO_OUTPUT_PATH)/goosedbversion.xo.go
	rm $(XO_OUTPUT_PATH)/xo_db.xo.go


