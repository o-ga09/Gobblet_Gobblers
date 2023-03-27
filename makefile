.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:	#	build docker image to deploy
	docker build -t taiti09/gotodo:${DOCKER_TAG} \
			--target deploy ./
	
build-local:
	docker compose build --no-cache

build-dev:
	docker-compose -f docker-compose_dev.yml build

up-dev:
	docker-compose -f docker-compose_dev.yml up -d

up:
	docker compose up

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

test:
	go test -race -suffle=on ./...

help:
	@grep -E '[^a-zA-Z_-]+:,*?## .*$$' $(MAKEFILE_LIST) | \
			awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

migrate:
	mysqldef -u todo -p P@ssw0rd -h 127.0.0.1 -P 33306 todo < ./_tools/mysql/schema.sql

dry-migrate:
	mysqldef -u todo -p P@ssw0rd -h 127.0.0.1 -P 33306 todo --dry-run < ./_tools/mysql/schema.sql

generate:
	go generate ./...