# --------------------------
# Env Management
# --------------------------
# Uses Docker
quickstart: update-env-file tidy db-docker-up sleep seed product-api-docker-up
# Uses Local GO
quickstart-local: update-env-file tidy db-docker-up sleep seed run

tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

update-env-file:
	@echo 'Updating .env from .env.example 🖋️...'
	@cp .env.example .env
	@echo '.env Updated ✨'

# --------------------------
# Run
# --------------------------
run ::
	@go run cmd/server/main.go

run-reload:
	@air -c .air.server-reloader.toml

product-api-docker-up ::
	docker compose up product-api

product-api-docker-down ::
	docker compose down product-api

product-api-docker-build:
	@echo 'Building images ...🛠️'
	@docker compose build product-api

db-docker-up ::
	docker compose up -d db-postgres

db-docker-down ::
	docker compose down db-postgres

docker-stop:
	docker compose down

docker-cleanup:
	docker compose down --remove-orphans

sleep:
	sleep 5

# --------------------------
# Test
# --------------------------
test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic
