.SILENT:

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/ocp-offer-api ./cmd/ocp-offer-api/main.go

# --

.PHONY: test
all:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
coverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'


.DEFAULT_GOAL := build
