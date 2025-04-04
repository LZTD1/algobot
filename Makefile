.PHONY: gen
gen:
	protoc --go_out=. \
	       --go_opt=paths=source_relative \
	       --go-grpc_out=. \
	       --go-grpc_opt=paths=source_relative \
	       ./protos/*.proto

.PHONY: dev
dev:
	go run ./cmd/algobot/main.go -config=./config/local.yaml


.PHONY: mock-gen
mock-gen:
	cd test && go generate ./...

.PHONY: grpc-gen
grpc-gen:
	protoc --go_out=. \
	       --go_opt=paths=source_relative \
	       --go-grpc_out=. \
	       --go-grpc_opt=paths=source_relative \
	       ./protos/*.proto

.PHONY: migrate
migrate:
	go run ./cmd/migrator/main.go -migrations-path=./migrations -storage-path=./storage/storage.db