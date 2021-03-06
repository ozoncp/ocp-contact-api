GOOS := $(shell go env GOOS)
ifeq ($(GOOS),windows)
	EXT := .exe
endif

.PHONY: deploy
deploy: vendor-proto .generate .build .compose-build .compose-up .migrate

.PHONY: start
start: .compose-build .compose-up .migrate

.PHONY: .compose-build
.compose-build:
	docker-compose build

.PHONY: .compose-up
.compose-up:
	docker compose up -d

.PHONY: stop
stop: .compose-stop

.PHONY: .compose-stop
.compose-stop:
	docker compose stop

.PHONY: .migrate
.migrate:
	 goose -dir ./migrations postgres "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable" up

.PHONY: build
build: vendor-proto .generate .build

.PHONY: generate
generate: .vendor-proto .generate

.PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-contact-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-contact-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-contact-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-contact-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-contact-api \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-contact-api/ocp-contact-api.proto
		mv pkg/ocp-contact-api/github.com/ozoncp/ocp-contact-api/pkg/ocp-contact-api/* pkg/ocp-contact-api/
		rm -rf pkg/ocp-contact-api/github.com
		mkdir -p cmd/ocp-contact-api

.PHONY: .build
.build:
		CGO_ENABLED=0 GOOS=$(GOOS) go build -o bin/ocp-contact-api$(EXT) cmd/ocp-contact-api/main.go

.PHONY: install
install: build .install

.PHONY: .install
install:
		go install cmd/ocp-contact-api/main.go

.PHONY: vendor-proto
vendor-proto: .vendor-proto

.PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/ocp-contact-api
		cp api/ocp-contact-api/ocp-contact-api.proto vendor.protogen/api/ocp-contact-api
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate ;\
		fi


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u google.golang.org/grpc
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u github.com/envoyproxy/protoc-gen-validate
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go install github.com/envoyproxy/protoc-gen-validate

.PHONY: docker
docker:
		docker build --no-cache -t ocp-contact-api:v1 .