# --------------------------
# Env Management
# --------------------------
quickstart: update-env-file tidy docker-up sleep seed run

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


test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down


sleep:
	sleep 5