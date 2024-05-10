.PHONY: comment/protoc
comment/protoc:
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

.PHONY: comment/new-migration
comment/new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: comment/migrateup
comment/migrateup:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54325/postgres?sslmode=disable" -verbose up

.PHONY: comment/migrateup1
comment/migrateup1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54325/postgres?sslmode=disable" -verbose up 1

.PHONY: comment/migratedown
comment/migratedown:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54325/postgres?sslmode=disable" -verbose down

.PHONY: comment/migratedown1
comment/migratedown1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54325/postgres?sslmode=disable" -verbose down 1

.PHONY: comment/wire
comment/wire:
	wire ./cmd

.PHONY: comment/mock
comment/mock:
	mockgen -source=./internal/repository/comment_repository.go -destination=./mock/mock_comment_repository.go -package=mock
	mockgen -source=./proto/comment_grpc.pb.go -destination=./mock/mock_comment_grpc.pb.go -package=mock

.PHONY: comment/test
comment/test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out

.PHONY: comment/build
comment/build:
	docker build -f Dockerfile -t techbranch/comment:latest ../..

.PHONY: comment/run
comment/run:
	go run ./cmd

