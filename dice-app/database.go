package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/XSAM/otelsql"
	_ "modernc.org/sqlite"

	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var db *sql.DB
var dbStatsReg metric.Registration

func initDB(ctx context.Context) error {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = "./data/chinook.db"
	}

	var err error
	db, err = otelsql.Open("sqlite", dbPath,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
	)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	dbStatsReg, err = otelsql.RegisterDBStatsMetrics(db,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
	)
	if err != nil {
		// Non-fatal: metrics might not be critical for startup
		return nil
	}
	return nil
}

func closeDB(ctx context.Context) error {
	if dbStatsReg != nil {
		_ = dbStatsReg.Unregister()
	}
	if db != nil {
		return db.Close()
	}
	return nil
}

func getDB() *sql.DB {
	return db
}
