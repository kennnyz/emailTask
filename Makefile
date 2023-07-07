createdb:
	docker run --name email_users -e POSTGRES_DB=email_users -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
migrateup:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/email_users?sslmode=disable" -verbose up

.PHONY: createdb
