createdb:
	docker run --name mailService -e POSTGRES_DB=user_mails -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
migrateup:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/user_mails?sslmode=disable" -verbose up

.PHONY: createdb migrateup
