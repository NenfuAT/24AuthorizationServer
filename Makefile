-include .env

build:
	docker compose build

up:
	docker compose up -d

log:
	docker compose logs

go:
	docker exec -it $(AUTHORIZATION_BACK_CONTAINER_HOST) /bin/sh

react:
	docker exec -it $(AUTHORIZATION_FRONT_CONTAINER_HOST) /bin/sh

db:
	docker exec -it ${MYSQL_CONTAINER_HOST} mysql -u ${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD} -D ${MYSQL_DATABASE}
