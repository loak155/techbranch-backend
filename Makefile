.PHONY: protoc
protoc:
	for file in $$(find proto -name '*.proto'); do \
		protoc \
		-I $$(dirname $$file) \
		-I ./third_party \
		--proto_path=. \
		--go_out=$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$$(dirname $$file)" --validate_opt=paths=source_relative \
		--grpc-gateway_out=$$(dirname $$file) --grpc-gateway_opt=paths=source_relative \
		$$file; \
	done

.PHONY: run-gateway
run-gateway:
	go run ./cmd/gateway/

.PHONY: run-user
run-user:
	go run ./cmd/user/

.PHONY: run-auth
run-auth:
	go run ./cmd/auth/

.PHONY: run-article
run-article:
	go run ./cmd/article/

.PHONY: run-bookmark
run-bookmark:
	go run ./cmd/bookmark/

.PHONY: run-comment
run-comment:
	go run ./cmd/comment/

.PHONY: new-migration
new-migration:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrateup
migrateup:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up 1

.PHONY: migratedown
migratedown:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose down 1

.PHONY: wire
wire:
	wire ./cmd/user
	wire ./cmd/auth
	wire ./cmd/article
	wire ./cmd/bookmark
	wire ./cmd/comment

.PHONY: build
build:
	docker build -f build/docker/user/Dockerfile -t techbranch/user .
	docker build -f build/docker/auth/Dockerfile -t techbranch/auth .
	docker build -f build/docker/article/Dockerfile -t techbranch/article .
	docker build -f build/docker/bookmark/Dockerfile -t techbranch/bookmark .
	docker build -f build/docker/comment/Dockerfile -t techbranch/comment .
	docker build -f build/docker/gateway/Dockerfile -t techbranch/gateway .

.PHONY: up
up:
	docker-compose up --build

.PHONY: down
down:
	docker-compose down

.PHONY: mock
mock:
	mockgen -source=./internal/user/repository/user_repository.go -destination=./mock/mock_user_repository.go -package=mock
	mockgen -source=./internal/article/repository/article_repository.go -destination=./mock/mock_article_repository.go -package=mock
	mockgen -source=./internal/comment/repository/comment_repository.go -destination=./mock/mock_comment_repository.go -package=mock
	mockgen -source=./internal/bookmark/repository/bookmark_repository.go -destination=./mock/mock_bookmark_repository.go -package=mock
	mockgen -source=./proto/user/user_grpc.pb.go -destination=./mock/mock_user_grpc.pb.go -package=mock
	mockgen -source=./proto/auth/auth_grpc.pb.go -destination=./mock/mock_auth_grpc.pb.go -package=mock
	mockgen -source=./proto/article/article_grpc.pb.go -destination=./mock/mock_article_grpc.pb.go -package=mock
	mockgen -source=./proto/comment/comment_grpc.pb.go -destination=./mock/mock_comment_grpc.pb.go -package=mock

.PHONY: test
test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out