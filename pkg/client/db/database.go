package db

import (
	"database/sql"
	"time"

	"gofiber-starterkit/pkg/utils"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func New() *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(utils.DatabaseConnectionString()),
	))

	pg := bun.NewDB(sqldb, pgdialect.New())
	pg.AddQueryHook(bundebug.NewQueryHook(bundebug.WithEnabled(true)))
	pg.SetMaxOpenConns(25)
	pg.SetMaxIdleConns(10)
	pg.SetConnMaxLifetime(5 * time.Minute)
	pg.SetConnMaxIdleTime(5 * time.Minute)

	if err := pg.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Debug().Msg("Connected to database successfully")
	return pg
}
