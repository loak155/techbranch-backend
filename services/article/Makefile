.PHONY: article/protoc
article/protoc:
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

.PHONY: article/new-migration
article/new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: article/migrateup
article/migrateup:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54323/postgres?sslmode=disable" -verbose up

.PHONY: article/migrateup1
article/migrateup1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54323/postgres?sslmode=disable" -verbose up 1

.PHONY: article/migratedown
article/migratedown:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54323/postgres?sslmode=disable" -verbose down

.PHONY: article/migratedown1
article/migratedown1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:54323/postgres?sslmode=disable" -verbose down 1

.PHONY: article/wire
article/wire:
	wire ./cmd

.PHONY: article/mock
article/mock:
	mockgen -source=./internal/repository/article_repository.go -destination=./mock/mock_article_repository.go -package=mock
	mockgen -source=./proto/article_grpc.pb.go -destination=./mock/mock_article_grpc.pb.go -package=mock

.PHONY: article/test
article/test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out

.PHONY: article/build
article/build:
	docker build -f Dockerfile -t techbranch/article:latest ../..

.PHONY: article/run
article/run:
	go run ./cmd

