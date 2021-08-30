package database

import (
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
)

func NewPostgres(dsn, driver string) *sqlx.DB {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create database connection")

		return nil
	}

	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msgf("failed ping the database")

		return nil
	}

	return db
}
