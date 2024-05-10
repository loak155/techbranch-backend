.PHONY: user/protoc
user/protoc:
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

.PHONY: user/new-migration
user/new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: user/migrateup
user/migrateup:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable" -verbose up

.PHONY: user/migrateup1
user/migrateup1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable" -verbose up 1

.PHONY: user/migratedown
user/migratedown:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable" -verbose down

.PHONY: user/migratedown1
user/migratedown1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable" -verbose down 1

.PHONY: user/wire
user/wire:
	wire ./cmd

.PHONY: user/mock
user/mock:
	mockgen -source=./internal/repository/user_repository.go -destination=./mock/mock_user_repository.go -package=mock
	mockgen -source=./proto/user_grpc.pb.go -destination=./mock/mock_user_grpc.pb.go -package=mock

.PHONY: user/test
user/test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out

.PHONY: user/build
user/build:
	docker build -f Dockerfile -t techbranch/user:latest ../..

.PHONY: user/run
user/run:
	go run ./cmd

