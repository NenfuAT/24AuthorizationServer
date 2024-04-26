-include .env

build:
	docker compose build

up:
	docker compose up -d

log:
	docker compose logs

go:
	docker exec -it $(AUTHORIZATION_CONTAINER_HOST) /bin/sh