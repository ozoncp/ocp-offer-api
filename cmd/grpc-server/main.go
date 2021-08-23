package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/server"
	"github.com/ozoncp/ocp-offer-api/internal/tracer"
)

var (
	batchSize uint = 2
)

func main() {
	log.Info().
		Str("version", cfg.Project.Version).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	db := createDB()
	defer db.Close()

	tracer.InitTracing("ocp_offer_api")

	if err := server.NewGrpcServer(db, batchSize).Start(); err != nil {
		log.Fatal().Err(err)
	}
}

func createDB() *sqlx.DB {
	dataSourceName := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Open(cfg.Database.Driver, dataSourceName)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create database connection")
		return nil
	}

	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msgf("failed ping the database")
		return nil
	}

	return db
}
