// Package db contains helper functions for working with db.
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
	"go.uber.org/zap"

	"github.com/dragonator/rental-service/pkg/config"
)

// OpenPGX opens a new DB connection using the pgx driver.
func OpenPGX(config *config.Config, logger *zap.Logger) (*sql.DB, error) {
	dsn := connectionString(config.Database)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}

	loggerAdapter := zapadapter.New(logger)
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)

	return db, nil
}

// Close closes a DB connection. Panics if closing fails.
func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func connectionString(dbConfig *config.Database) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
}
