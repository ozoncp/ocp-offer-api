package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/database"
	"github.com/ozoncp/ocp-offer-api/internal/server"
	"github.com/ozoncp/ocp-offer-api/internal/tracer"
	"github.com/pressly/goose/v3"
)

var (
	batchSize uint = 2
)

func main() {
	migration := flag.String("migration", "", "Defines the migration start option")
	flag.Parse()

	log.Info().
		Str("version", cfg.Project.Version).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db := database.NewPostgres(dsn, cfg.Database.Driver)

	if *migration != "" {
		migrate(db.DB, *migration)
	}

	tracer.InitTracing("ocp_offer_api")

	if err := server.NewGrpcServer(db, batchSize).Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed creating gRPC server")
	}

	db.Close()
}

func migrate(db *sql.DB, command string) {
	switch command {
	case "up":
		if err := goose.Up(db, "migrations"); err != nil {
			log.Fatal().Err(err).Msg("Migration failed")
		}
	case "down":
		if err := goose.Down(db, "migrations"); err != nil {
			log.Fatal().Err(err).Msg("Migration failed")
		}

	default:
		log.Warn().Msgf("Invalid command for 'migration' flag: '%v'", command)
	}
}
