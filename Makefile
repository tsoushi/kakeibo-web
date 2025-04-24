dev/run/server:
	. ./.env
	docker compose up -d

dev/run/down:
	docker compose down

restart: dev/run/down dev/run/server
