package postgres

import (
	"context"
	"fmt"

	"github.com/grum261/beer/internal/models"
	"github.com/pkg/errors"
)

const (
	userInsert = `
	INSERT INTO users (username, email, password_hash, bio, avatart) VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	userSelectPassword = `
	SELECT id, password_hash FROM users
	WHERE %s = $1`
)

func (q *Queries) InsertUser(ctx context.Context, p models.UserCreateParams) (int, error) {
	var userID int

	if err := q.db.QueryRow(
		ctx, userInsert, p.Username, p.Email,
		p.Password, p.Bio, p.Avatar,
	).Scan(&userID); err != nil {
		return userID, errors.Wrap(err, "postgres.InsertUser")
	}

	return userID, nil
}

func (q *Queries) SelectPassword(ctx context.Context, arg, field string) (int, string, error) {
	var sql string

	switch field {
	case "username", "email":
		sql = fmt.Sprintf(userSelectPassword, field)
	default:
		return 0, "", errors.New("postgres.SelectPassword: invalid field for password selection")
	}

	var (
		hashedPass string
		userID     int
	)

	if err := q.db.QueryRow(ctx, sql, arg).Scan(&userID, &hashedPass); err != nil {
		return userID, hashedPass, errors.Wrap(err, "postgres.SelectPassword")
	}

	return userID, hashedPass, nil
}
