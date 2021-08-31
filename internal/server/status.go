package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/rs/zerolog/log"
)

func createStatusServer(addr string, isReady *atomic.Value) *http.Server {
	mux := http.DefaultServeMux

	mux.HandleFunc(cfg.Status.LivenessPath, livenessHandler)
	mux.HandleFunc(cfg.Status.ReadinessPath, readinessHandler(isReady))
	mux.HandleFunc(cfg.Status.VersionPath, versionHandler)

	statusServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return statusServer
}

func livenessHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func versionHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]interface{}{
		"name":        cfg.Project.Name,
		"debug":       cfg.Project.Debug,
		"environment": cfg.Project.Environment,
		"version":     cfg.Project.Version,
		"commitHash":  cfg.Project.CommitHash,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("Service information encoding error")
	}
}
