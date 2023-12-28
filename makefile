create_migration:
	# creates a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

postgres_up:
	# start postgres container.
	docker compose up -d 

postgres_down:
	# stop postgres container
	docker compose down 

db_up:
	docker exec -it pg_local createdb --username=root --owner=root fingreat

db_down:
	docker exec -it pg_local dropdb --username=root fingreat

migrate_up:
	# Make migration
	migrate -path db/migrations -database "postgres://<user>:<password>@localhost:5432/fingreat?sslmode=disable" up

migrate_down:
	# rollback migration
	migrate -path db/migrations -database "postgres://<user>:<password>@localhost:5432/fingreat?sslmode=disable" down $(count)

download_deps:
	# download dependencies
	go mod download


sqlc:
  # generate go sql
	sqlc generate 

start_prod:	
	# start server in prod
	./fingreat

start_dev:
	# start server
	CompileDaemon -command="./fingreat"

test:
  # run covrage test.[fix this]
	go test -run github/tdadadavid/fingreat/backend/db/tests