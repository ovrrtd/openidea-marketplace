.PHONY: runServer
runServer:
	cd cmd && go build -o main && ./main

.PHONY: migrateUp
migrateUp:
	migrate -database "postgres://root:root@localhost:5432/ecomm?sslmode=disable" -path db/migrations up

.PHONY: migrateDown
migrateDown:
	migrate -database "postgres://root:root@localhost:5432/ecomm?sslmode=disable" -path db/migrations down