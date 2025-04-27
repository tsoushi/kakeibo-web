include .env

local/up:
	. ./.env
	docker compose -f docker-compose-local.yml up -d

local/down:
	docker compose -f docker-compose-local.yml down

local/restart: local/down local/up

server_gqlgen:
	make -C ./server gen

install_tools:
	mkdir ./bin
	wget -O - https://github.com/sqldef/sqldef/releases/latest/download/mysqldef_linux_amd64.tar.gz | tar -xvz -C ./bin

local/create_table:
	./bin/mysqldef -u ${MYSQL_ROOT_USER} -p ${MYSQL_ROOT_PASSWORD} -h localhost -P ${MYSQL_PORT} ${MYSQL_DATABASE} < ./server/schema.sql