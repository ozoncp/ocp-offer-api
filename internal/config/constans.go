package config

const (
	// GRPC environment constants
	GRPC_HOST          = "GRPC_HOST"
	GRPC_PORT          = "GRPC_PORT"
	GRPC_MAX_CONN_IDLE = "GRPC_MAX_CONN_IDLE"
	GRPC_TIMEOUT       = "GRPC_TIMEOUT"
	GRPC_CONN_AGE      = "GRPC_CONN_AGE"

	// Gataway environment constants
	GATEWAY_HOST = "GATEWAY_HOST"
	GATEWAY_PORT = "GATEWAY_PORT"

	// Metrics environment constants
	METRICS_HOST = "METRICS_HOST"
	METRICS_PORT = "METRICS_PORT"
	METRICS_PATH = "METRICS_PATH"

	// DATABASE environment constants
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_NAME     = "DATABASE_NAME"
	DATABASE_SSL_MODE = "DATABASE_SSL_MODE"
	DATABASE_DRIVER   = "DATABASE_DRIVER"
)
