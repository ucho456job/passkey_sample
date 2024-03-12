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
	@migrate -path ./db/migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" up

migrate_down:
	@migrate -path ./db/migrations -database "postgres://user:password@localhost:5432/passkey_db?sslmode=disable" down

insert_data:
	@PGPASSWORD=password psql -U user -h localhost -d passkey_db -f ./db/data/insert.sql

clear_data:
	@PGPASSWORD=password psql -U user -h localhost -d passkey_db -f ./db/data/clear.sql

migrate_reset:
	$(MAKE) migrate_down
	$(MAKE) migrate_up