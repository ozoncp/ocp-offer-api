LOCAL_BIN:=$(CURDIR)/bin


run:
	go run cmd/ocp-offer-api/main.go


lint:
	golint ./...


.PHONY: test
test: all

all:
	go test -v -race -timeout 30s -coverprofile cover.out ./...

coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'


.PHONY: generate
generate: .vendor-proto .generate

.PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/ocp-offer-api
		cp api/ocp-offer-api/*.proto vendor.protogen/api/ocp-offer-api
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/github.com/envoyproxy ]; then \
			mkdir -p vendor.protogen/github.com/envoyproxy &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/github.com/envoyproxy/protoc-gen-validate;\
		fi


.PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-offer-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-offer-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-offer-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-offer-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--validate_out lang=go:pkg/ocp-offer-api \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-offer-api/ocp-offer-api.proto
		mv pkg/ocp-offer-api/github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api/* pkg/ocp-offer-api/
		rm -rf pkg/ocp-offer-api/github.com
		mkdir -p cmd/ocp-offer-api
		cd pkg/ocp-offer-api && ls go.mod || go mod init github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api && go mod tidy


.PHONY: build
build: generate .build


.PHONY: .build
.build:
		go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/ocp-offer-api ./cmd/ocp-offer-api/main.go


.PHONY: install
install: build .install

.PHONY: .install
install:
		go install cmd/grpc-server/main.go


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init github.com/ozoncp/ocp-offer-api
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u google.golang.org/grpc
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u github.com/envoyproxy/protoc-gen-validate
		go get -u -v github.com/yoheimuta/protolint/cmd/protolint
		go get -u google.golang.org/grpc/test/bufconn
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go install github.com/envoyproxy/protoc-gen-validate


.DEFAULT_GOAL := build
