package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"sockets-tc/internal/adapters/repository/postgres/db"
	"sockets-tc/internal/ports"
)

var (
	_ ports.EventRepository       = (*Repository)(nil)
	_ ports.ReservationRepository = (*Repository)(nil)
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *Repository) withTx(ctx context.Context, fn func(q *db.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(r.queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
