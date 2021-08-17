module github.com/ozoncp/ocp-offer-api

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.40.0
)

replace github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api => ./pkg/ocp-offer-api
