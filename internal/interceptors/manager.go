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
		Name: "grpc_microservice_requests_total",
		Help: "The total number of incoming gRPC requests",
	})
)

// InterceptorManager struct.
type InterceptorManager struct {
}

// NewInterceptorManager InterceptorManager constructor.
func NewInterceptorManager() *InterceptorManager {
	return &InterceptorManager{}
}

// Logger Interceptor.
func (im *InterceptorManager) Logger(
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
