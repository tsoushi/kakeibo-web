include .env

local/up:
	docker compose -f docker-compose-local.yml up -d

local/down:
	docker compose -f docker-compose-local.yml down

local/restart: local/down local/up

server_gqlgen:
	make -C ./server gen

install_tools/linux:
	mkdir ./bin
	wget -O - https://github.com/sqldef/sqldef/releases/latest/download/mysqldef_licux_amd64.tar.gz | tar -xvz -C ./bin

install_tools/mac:
	brew install sqldef/sqldef/mysqldef

local/create_table:
	mysqldef -u ${MYSQL_ROOT_USER} -p ${MYSQL_ROOT_PASSWORD} -h localhost -P ${MYSQL_PORT} ${MYSQL_DATABASE} < ./server/migrate/schema.sql
	mysql -u ${MYSQL_ROOT_USER} -p${MYSQL_ROOT_PASSWORD} -h 127.0.0.1 -P ${MYSQL_PORT} ${MYSQL_DATABASE} < ./server/migrate/seed.sql