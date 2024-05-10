.PHONY: auth/protoc
auth/protoc:
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

.PHONY: auth/wire
auth/wire:
	wire ./cmd

.PHONY: auth/mock
auth/mock:
	mockgen -source=./proto/auth_grpc.pb.go -destination=./mock/mock_auth_grpc.pb.go -package=mock

.PHONY: auth/test
auth/test:
	go test ./... -v -cover -covermode=set -coverprofile=coverage.out

.PHONY: auth/build
auth/build:
	docker build -f Dockerfile -t techbranch/auth:latest ../..

.PHONY: auth/run
auth/run:
	go run ./cmd

