DB_URL=postgresql://root:secret@localhost:5432/techbranch?sslmode=disable

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrateup
migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

.PHONY: migratedown
migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

.PHONY: protoc
protoc:
	rm -f pkg/pb/*.go
	rm -f docs/swagger/techbranch.swagger.json
	protoc \
	-I third_party \
	--proto_path=api/proto \
	--go_out=pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--validate_out="lang=go:pkg/pb" --validate_opt=paths=source_relative \
	--grpc-gateway_out=allow_delete_body=true:pkg/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=techbranch \
	api/proto/*.proto
	statik -src=./docs/swagger -dest=./docs/swagger

.PHONY: mockgen
mockgen:
	mockgen -source=./internal/repository/article_repository.go -destination=./mock/mock_article_repository.go -package=mock

.PHONY: test
test:
	go test -v -cover -short ./...

.PHONY: run
run:
	ENV=local go run ./cmd

.PHONY: dbdocs
db_docs:
	dbdocs build docs/db/db.dbml

.PHONY: dbml2sql
dbml2sql:
	dbml2sql --postgres -o docs/db/schema.sql docs/db/db.dbml