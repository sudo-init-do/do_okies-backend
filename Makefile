run:
	go run ./apps/api

dc-up:
	docker compose -f infra/docker/docker-compose.yml up --build

dc-down:
	docker compose -f infra/docker/docker-compose.yml down
