package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_total_requests",
		Help: "The total number of incoming HTTP requests",
	})
)

func createGatewayServer(grpcAddr, gatewayAddr string) *http.Server {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal().Msgf("failed to dial server: %w", err)
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterOcpOfferApiServiceHandler(context.Background(), mux, conn); err != nil {
		log.Fatal().Err(err)
	}

	gatewayServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}

			httpTotalRequests.Inc()

			mux.ServeHTTP(w, r)
		}),
	}

	return gatewayServer
}
