package main

import (
	"context"
	"net/http"

	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	api "github.com/ozoncp/ocp-offer-api/internal/api"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	desc "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
)

const (
	grpcPort           = ":9090"
	grpcServerEndpoint = "localhost:9090"
)

func runGrpcServer() error {
	l, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	defer l.Close()

	grpcServer := grpc.NewServer()

	desc.RegisterOcpOfferApiServiceServer(grpcServer, api.NewOfferAPI(repo.NewRepo()))

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}

	return nil
}

func runHttpServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpOfferApiServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	go runHttpServer()

	if err := runGrpcServer(); err != nil {
		log.Fatal().Err(err)
	}
}
