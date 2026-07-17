package repository

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGRepository struct {
	mu      sync.Mutex
	pgxPool *pgxpool.Pool
}

func NewPGRepository(connectStr string) (*PGRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), connectStr)
	if err != nil {
		return nil, err
	}

	return &PGRepository{mu: sync.Mutex{}, pgxPool: pool}, nil
}
