package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ozoncp/ocp-offer-api/internal/api"
	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/interceptors"
	"github.com/ozoncp/ocp-offer-api/internal/producer"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	pb "github.com/ozoncp/ocp-offer-api/pkg/ocp-offer-api"
)

type GrpcServer struct {
	db        *sqlx.DB
	batchSize uint
}

func NewGrpcServer(db *sqlx.DB, batchSize uint) *GrpcServer {
	return &GrpcServer{
		db:        db,
		batchSize: batchSize,
	}
}

func (s *GrpcServer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	im := interceptors.NewInterceptorManager()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.Gateway.Host, cfg.Gateway.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.GRPC.Host, cfg.GRPC.Port)
	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)

	gatewayServer := createGatewayServer(grpcAddr, gatewayAddr)

	go func() {
		log.Info().Msgf("Gateway server is running on %s", gatewayAddr)
		if err := gatewayServer.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Failed running gateway server")
			cancel()
		}
	}()

	metricsServer := createMetricsServer(metricsAddr)

	go func() {
		log.Info().Msgf("Metrics server is running on %s", metricsAddr)
		if err := metricsServer.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Failed running metrics server")
			cancel()
		}
	}()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.GRPC.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.GRPC.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.GRPC.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.GRPC.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			im.Logger,
		)),
	)

	repo := repo.NewRepo(s.db, s.batchSize)
	dataProducer, err := producer.New(ctx, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Capacity)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create a producer")
	}

	pb.RegisterOcpOfferApiServiceServer(grpcServer, api.NewOfferAPI(repo, dataProducer))
	grpc_prometheus.Register(grpcServer)

	go func() {
		log.Info().Msgf("GRPC Server is listening on: %s", grpcAddr)
		if err := grpcServer.Serve(l); err != nil {
			log.Fatal().Err(err).Msg("Failed running gRPC server")
		}
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info().Msgf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Info().Msgf("ctx.Done: %v", done)
	}

	if err := gatewayServer.Shutdown(ctx); err != nil {
		log.Info().Msgf("gatewayServer.Shutdown: %v", err)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Info().Msgf("metricsServer.Shutdown: %v", err)
	}

	grpcServer.GracefulStop()
	log.Info().Msgf("Server shut down correctly")

	return nil
}