// Package postgres ...
package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/n-r-w/lg"
	"github.com/n-r-w/nerr"
)

const (
	defaultMaxConns           = 10
	defaultConnAttempts       = 10
	defaultMaxMaxConnIdleTime = time.Second
	defaultConnTimeout        = time.Second * 5
	defaultStatementTimeout   = time.Second * 5
	defaultReconnectTimeout   = time.Second
)

type Service struct {
	maxConns           int
	connAttempts       int
	maxMaxConnIdleTime time.Duration
	connTimeout        time.Duration
	statementTimeout   time.Duration
	reconnectTimeout   time.Duration

	Pool *pgxpool.Pool
}

// Url - алиас для корректной работы google wire
type Url string

func New(url Url, logger lg.Logger, options ...Option) (*Service, error) {
	pg := &Service{
		maxConns:           defaultMaxConns,
		connAttempts:       defaultConnAttempts,
		maxMaxConnIdleTime: defaultMaxMaxConnIdleTime,
		connTimeout:        defaultConnTimeout,
		statementTimeout:   defaultStatementTimeout,
		reconnectTimeout:   defaultReconnectTimeout,
	}

	for _, opt := range options {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(string(url))
	if err != nil {
		return nil, nerr.New("pgxpool parse config error", err)
	}

	poolConfig.MaxConnIdleTime = pg.maxMaxConnIdleTime
	poolConfig.MaxConns = int32(pg.maxConns)

	logger.Info("connecting to database...")
	connAttempts := pg.connAttempts
	for connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		logger.Warn("connecting to database, attempts left: %d", connAttempts)

		time.Sleep(pg.reconnectTimeout)

		connAttempts--
	}

	if err != nil {
		return nil, nerr.New("database connection error", err)
	}

	logger.Info("database connection ok")

	return pg, nil
}

func (p *Service) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
