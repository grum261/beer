package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type pgdb interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Queries struct {
	db          pgdb
	stmtBuilder sq.StatementBuilderType
}

func NewQueries(db pgdb) *Queries {
	return &Queries{
		db:          db,
		stmtBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db:          tx,
		stmtBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
