package config

const (
	// GRPC environment constants.
	GrpcHost        = "GRPC_HOST"
	GrpcPort        = "GRPC_PORT"
	GrpcMaxConnIdle = "GRPC_MAX_CONN_IDLE"
	GrpcTimeout     = "GRPC_TIMEOUT"
	GrpcConnAge     = "GRPC_CONN_AGE"

	// Gataway environment constants.
	GatewayHost = "GATEWAY_HOST"
	GatewayPort = "GATEWAY_PORT"

	// Metrics environment constants.
	MetricsHost = "METRICS_HOST"
	MetricsPort = "METRICS_PORT"
	MetricsPath = "METRICS_PATH"

	// DATABASE environment constants.
	DatabaseHost     = "DATABASE_HOST"
	DatabasePort     = "DATABASE_PORT"
	DatabaseUser     = "DATABASE_USER"
	DatabasePassword = "DATABASE_PASSWORD"
	DatabaseName     = "DATABASE_NAME"
	DatabaseSslMode  = "DATABASE_SSL_MODE"
	DatabaseDriver   = "DATABASE_DRIVER"
)
