package repository

import (
	"context"

	"github.com/grum261/beer/internal/models"
	"github.com/grum261/beer/internal/repository/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Brewery struct {
	q *postgres.BreweryQueries
}

func NewBrewery(pool *pgxpool.Pool) *Brewery {
	return &Brewery{
		q: &postgres.BreweryQueries{
			Queries: postgres.NewQueries(pool),
		},
	}
}

func (b *Brewery) CreateBrewery(ctx context.Context, p models.BreweryCreateParams) (int, error) {
	return b.q.InsertBrewery(ctx, p)
}
