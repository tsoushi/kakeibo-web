local/up:
	. ./.env
	docker compose -f docker-compose-local.yml up -d

local/down:
	docker compose -f docker-compose-local.yml down

local/restart: local/down local/up
