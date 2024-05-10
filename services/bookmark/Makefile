.PHONY: bookmark/protoc
bookmark/protoc:
	for file in $$(find proto -name '*.proto'); do \
		protoc \
		-I $$(dirname $$file) \
		-I ../../third_party \
		--proto_path=. \
		--go_out=$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$$(dirname $$file)" --validate_opt=paths=source_relative \
		--grpc-gateway_out=allow_delete_body=true:$$(dirname $$file) --grpc-gateway_opt=paths=source_relative \
		$$file; \
	done

.PHONY: bookmark/new-migration
bookmark/new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: bookmark/migrateup
bookmark/migrateup:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54324/postgres?sslmode=disable" -verbose up

.PHONY: bookmark/migrateup1
bookmark/migrateup1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54324/postgres?sslmode=disable" -verbose up 1

.PHONY: bookmark/migratedown
bookmark/migratedown:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54324/postgres?sslmode=disable" -verbose down

.PHONY: bookmark/migratedown1
bookmark/migratedown1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54324/postgres?sslmode=disable" -verbose down 1

.PHONY: bookmark/wire
bookmark/wire:
	wire ./cmd

.PHONY: bookmark/mock
bookmark/mock:
	mockgen -source=./internal/repository/bookmark_repository.go -destination=./mock/mock_bookmark_repository.go -package=mock

.PHONY: bookmark/test
bookmark/test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out

.PHONY: bookmark/build
bookmark/build:
	docker build -f Dockerfile -t techbranch/bookmark:latest ../..

.PHONY: bookmark/run
bookmark/run:
	go run ./cmd

