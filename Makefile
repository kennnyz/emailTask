createdb:
	docker run --name mailServiceDB -e POSTGRES_DB=user_emails -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
migrateup:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/user_emails?sslmode=disable" -verbose up
docker-up:
	docker-compose up -d

.PHONY: createdb migrateup
