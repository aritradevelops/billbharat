package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Connect() error
	Disconnect() error
	Health() error
	Tx() (*pgxpool.Pool, error)
}

func NotInitializedErr(what string) error {
	return fmt.Errorf("%s is not initialized, have you forgot to call Connect() ?", what)
}

type Postgres struct {
	connString string
	timeout    time.Duration
	// conn *pgx.Conn
	pool *pgxpool.Pool
}

func NewPostgres(connString string, timeout time.Duration) *Postgres {
	return &Postgres{
		connString: connString,
		timeout:    timeout,
	}
}

func (p *Postgres) Connect() error {
	connectionCtx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	// conn, err := pgx.Connect(connectionCtx, p.uri)
	pool, err := pgxpool.New(connectionCtx, p.connString)
	if err != nil {
		return fmt.Errorf("failed to connect to the database : %v", err)
	}
	p.pool = pool
	return nil
}
func (p *Postgres) Disconnect() error {
	if p.pool == nil {
		return NotInitializedErr("Postgres")
	}
	p.pool.Close()
	p.pool = nil
	return nil
}
func (p *Postgres) Health() error {
	if p.pool == nil {
		return NotInitializedErr("Postgres")
	}
	pingCtx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	err := p.pool.Ping(pingCtx)
	if err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}
	return nil
}

func (p *Postgres) Tx() (*pgxpool.Pool, error) {
	if p.pool == nil {
		return nil, NotInitializedErr("Postgres")
	}
	return p.pool, nil
}
