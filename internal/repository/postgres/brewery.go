package postgres

import (
	"context"

	"github.com/grum261/beer/internal/models"
	"github.com/pkg/errors"
)

type BreweryQueries struct {
	*Queries
}

func (q *BreweryQueries) InsertBrewery(ctx context.Context, p models.BreweryCreateParams) (int, error) {
	query, args, err := q.stmtBuilder.
		Insert("breweries").
		Columns(
			"name",
			"description",
			"founded_at",
		).
		Values(p.Name, p.Description, p.FoundedAt).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "postgres.InsertBrewery")
	}

	var breweryID int

	if err := q.db.QueryRow(ctx, query, args...).Scan(&breweryID); err != nil {
		return breweryID, errors.Wrap(err, "postgres.InsertBrewery")
	}

	return breweryID, nil
}
