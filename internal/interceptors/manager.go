package interceptors

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_microservice_total_requests",
		Help: "The total number of incoming gRPC requests",
	})
)

// InterceptorManager struct
type interceptorManager struct {
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager() *interceptorManager {
	return &interceptorManager{}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	totalRequests.Inc()
	start := time.Now()

	reply, err := handler(ctx, req)
	log.Info().
		Str("Method", info.FullMethod).
		Dur("latency", time.Since(start)).
		Send()

	return reply, err
}
