.PHONY: up down logs rebuild api dbshell

up:
\tdocker compose up -d

down:
\tdocker compose down

rebuild:
\tdocker compose build api && docker compose up -d api

logs:
\tdocker compose logs -f api

dbshell:
\tdocker compose exec -it postgres psql -U okies -d okiesdb
