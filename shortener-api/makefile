dependencies:
	go mod download

dep-install:
	go mod tidy
	# go mod vendor

migrate-up:
	migrate -path db/migration/ -database "postgres://$(SHORTENER_API_DB_USER):$(SHORTENER_API_DB_PASSWORD)@$(SHORTENER_API_DB_HOST):$(SHORTENER_API_DB_PORT)/$(SHORTENER_API_DB_NAME)?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration/ -database "postgres://$(SHORTENER_API_DB_USER):$(SHORTENER_API_DB_PASSWORD)@$(SHORTENER_API_DB_HOST):$(SHORTENER_API_DB_PORT)/$(SHORTENER_API_DB_NAME)?sslmode=disable" -verbose down