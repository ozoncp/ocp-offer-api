package server

import (
	"net/http"

	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func createMetricsServer(addr string) *http.Server {

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return metricsServer
}
