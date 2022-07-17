package postgres

import "time"

type Option func(*Service)

func MaxConns(n int) Option {
	return func(p *Service) {
		p.maxConns = n
	}
}

func ConnAttempts(n int) Option {
	return func(p *Service) {
		p.connAttempts = n
	}
}

func MaxMaxConnIdleTime(n time.Duration) Option {
	return func(p *Service) {
		p.maxMaxConnIdleTime = n
	}
}

func ConnTimeout(n time.Duration) Option {
	return func(p *Service) {
		p.connTimeout = n
	}
}

func StatementTimeout(n time.Duration) Option {
	return func(p *Service) {
		p.statementTimeout = n
	}
}

func ReconnectTimeout(n time.Duration) Option {
	return func(p *Service) {
		p.reconnectTimeout = n
	}
}
