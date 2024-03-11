dependencies_start:
	docker compose up -d

dependencies_stop:
	docker compose down

front_build:
	@cd template && npm run build

run: ## Run the project
	go run main.go

migrate_up: ## Run the migrations for local
	@migrate -path ./migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" up

migrate_down: ## Rollback the migrations for local
	@migrate -path ./migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" down

migrate_reset: ## Reset the migrations
	$(MAKE) migrate_down
	$(MAKE) migrate_up