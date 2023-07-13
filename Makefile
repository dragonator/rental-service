envfile ?= .env

-include $(envfile)

define export_envfile
ifneq ("$(wildcard $(1))","")
	export $(shell sed 's/=.*//' $(envfile))
endif
endef

$(eval $(call export_envfile,$(envfile)))

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