package postgres

import (
	"context"
	"database/sql"

	"be/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	gorm *gorm.DB
	sql  *sql.DB
}

func Connect(cfg config.PostgresConfig) (*DB, error) {
	g, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s, err := g.DB()
	if err != nil {
		return nil, err
	}
	return &DB{gorm: g, sql: s}, nil
}

func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return nil
	}
	return d.sql.PingContext(ctx)
}

func (d *DB) Gorm() *gorm.DB {
	return d.gorm
}
