package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"time"
)

// InitDbConn подключение к БД
func InitDbConn(ctx context.Context) (dbpool *pgxpool.Pool, err error) {
	url := "postgres://postgres:password@localhost:5432/ekaterina?sslmode=disable"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}
	// конфигурация
	cfg.MaxConns = int32(5)
	cfg.MinConns = int32(1)
	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.MaxConnLifetime = 24 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	// создаем пул
	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}
	return
}
