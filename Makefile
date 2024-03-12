dependencies_start:
	docker compose up -d

dependencies_stop:
	docker compose down

front_build:
	@cd template && npm install
	@cd template && npm run build

run:
	go run main.go

migrate_up:
	@migrate -path ./migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" up

migrate_down:
	@migrate -path ./migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" down

migrate_reset:
	$(MAKE) migrate_down
	$(MAKE) migrate_up