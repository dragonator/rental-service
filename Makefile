envfile ?= .env

-include $(envfile)

define export_envfile
    export $(shell sed 's/=.*//' $(1))
endef

ifneq ("$(wildcard $(envfile))", "")
    $(eval $(call export_envfile,$(envfile)))
endif

.PHONY: init
init:
	@cp .env.dist .env


.PHONY: db-up
db-up:
	@docker-compose up -d --no-build --remove-orphans postgres
	@docker-compose ps


.PHONY: volumes-down
volumes-down:
	@docker-compose down -v --remove-orphans
	@docker-compose ps


.PHONY: db-down
db-down:
	@docker-compose down --remove-orphans
	@docker-compose ps


.PHONY: generate
generate:
	@find . -name "*_mock_test.go" | xargs -r rm
	@go generate ./...


.PHONY: test
test: generate
	@go test ./...


.PHONY: server-start
server-start: generate
	@go run cmd/server/main.go