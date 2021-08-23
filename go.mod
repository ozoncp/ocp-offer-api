module github.com/ozoncp/ocp-offer-api

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/Shopify/sarama v1.29.1
	github.com/envoyproxy/protoc-gen-validate v0.6.1 // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/gertd/go-pluralize v0.1.7 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/go-hclog v0.16.2 // indirect
	github.com/hashicorp/go-plugin v1.4.2 // indirect
	github.com/hashicorp/yamux v0.0.0-20210707203944-259a57b3608c // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/labstack/echo/v4 v4.5.0
	github.com/lib/pq v1.10.2
	github.com/lyft/protoc-gen-star v0.5.3 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api v0.0.0-00010101000000-000000000000
	github.com/pressly/goose/v3 v3.0.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/spf13/afero v1.6.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/yoheimuta/protolint v0.32.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210820121016-41cdb8703e55 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210821163610-241b8fcbd6c8 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api => ./pkg/ocp-offer-api
