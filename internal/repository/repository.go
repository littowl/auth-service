package db

import (
	"auth-service/utils"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewDB(pool *pgxpool.Pool, logger *zap.Logger) *DB {
	db := &DB{
		pool:   pool,
		logger: logger,
	}

	return db
}

func DbStart(conn *utils.Config, logger *zap.Logger) *pgxpool.Pool {
	connStr := "host=" + conn.PgHost + " port=" + conn.PgPort + " user=" + conn.PgUser + " password=" + conn.PgPass + " dbname=" + conn.PgDB + " sslmode=disable"
	logger.Info("Service init: connecting to database...")
	urlExample := string(connStr)
	dbpool, err := pgxpool.Connect(context.Background(), urlExample)
	if err != nil {
		logger.Error("db pool connect error", zap.Error(err))
		return nil
	}
	return dbpool
}
