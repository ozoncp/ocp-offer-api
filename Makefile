GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

PGV_VERSION:="v0.6.1"
GOOGLEAPIS_VERSION="master"
BUF_VERSION:="v0.51.0"
GOBIN?=$(GOPATH)/bin

.PHONY: run
grpc-server:
	go run cmd/grpc-server/main.go

kafka-consumer:
	go run cmd/kafka-consumer/main.go


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $3}'


# ----------------------------------------------------------------

.PHONY: generate
generate: .vendor-proto .generate

.generate:
	@command -v buf 2>&1 > /dev/null || (mkdir -p $(GOBIN) && curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(shell uname -s)-$(shell uname -m) -o $(GOBIN)/buf && chmod +x $(GOBIN)/buf)
	$(eval PROTOS:=$(CURDIR)/protos)
	@[ -f $(PROTOS)/buf.yaml ] || PATH=$(GOBIN):$(PATH) buf mod init --doc -o $(PROTOS)
	PATH=$(GOBIN):$(PATH) buf generate $(PROTOS)

.vendor-proto:
	$(eval VENDOR:=$(CURDIR)/vendor.protogen)
	@[ -f $(VENDOR)/validate/validate.proto ] || (mkdir -p $(VENDOR)/validate/ && curl -sSL0 https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/$(PGV_VERSION)/validate/validate.proto -o $(VENDOR)/validate/validate.proto)
	@[ -f $(VENDOR)/google/api/http.proto ] || (mkdir -p $(VENDOR)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/http.proto -o $(VENDOR)/google/api/http.proto)
	@[ -f $(VENDOR)/google/api/annotations.proto ] || (mkdir -p $(VENDOR)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/annotations.proto -o $(VENDOR)/google/api/annotations.proto)
	cd pkg/ocp-offer-api && ls go.mod || go mod init github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api && go mod tidy

# ----------------------------------------------------------------

.PHONY: deps
deps: .deps .bin-deps

.deps:
	@[ -f go.mod ] || go mod init github.com/ozoncp/ocp-offer-api
	find . -name go.mod -exec bash -c 'pushd "$${1%go.mod}" && go mod tidy && popd' _ {} \;

.bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)

.PHONY: build
build: generate
	go mod download && CGO_ENABLED=0 GOOS=linux go build -tags='no_mysql no_sqlite3' -o ./bin/grpc-server ./cmd/grpc-server/main.go
	go mod download && CGO_ENABLED=0 GOOS=linux go build -tags='no_mysql no_sqlite3' -o ./bin/kafka-consumer ./cmd/kafka-consumer/main.go
