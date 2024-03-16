# to run server
.PHONY: runServer
runServer:
	cd cmd && go build -o main && ./main

# to run migration file
.PHONY: migrateUp
migrateUp:
	migrate -database "postgres://root:root@localhost:5432/ecomm?sslmode=disable" -path db/migrations up

# to run rollback migration
.PHONY: migrateDown
migrateDown:
	migrate -database "postgres://root:root@localhost:5432/ecomm?sslmode=disable" -path db/migrations down