package repository

import (
	"context"

	"github.com/pkg/errors"

	"github.com/grum261/beer/internal/models"
	"github.com/grum261/beer/internal/repository/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	q *postgres.UserQueries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		q: &postgres.UserQueries{
			Queries: postgres.NewQueries(pool),
		},
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, p models.UserCreateParams) (int, error) {
	userID, err := u.q.InsertUser(ctx, p)
	if err != nil {
		return userID, errors.Wrap(err, "repository.CreateUser")
	}

	return userID, nil
}

func (u *UserRepository) GetUserPasswordByName(ctx context.Context, username string) (int, string, error) {
	userID, hashedPass, err := u.q.SelectPassword(ctx, username, "username")
	if err != nil {
		return userID, hashedPass, errors.Wrap(err, "repository.GetUserPasswordByName")
	}

	return userID, hashedPass, nil
}

func (u *UserRepository) GetUserPasswordByEmail(ctx context.Context, email string) (int, string, error) {
	userID, hashedPass, err := u.q.SelectPassword(ctx, email, "email")
	if err != nil {
		return userID, hashedPass, errors.Wrap(err, "repository.GetUserPasswordByEmail")
	}

	return userID, hashedPass, nil
}

func (u *UserRepository) UpdateFriendStatus(
	ctx context.Context,
	senderID, receiverID int,
	status models.RequestStatus,
) error {
	if err := status.Validate(); err != nil {
		return errors.Wrap(err, "service.UpdateFriendRequest")
	}

	// TODO: rework to inert on conflict do update
	if status == models.StatusSent {
		return u.q.InsertFriendRequest(ctx, senderID, receiverID)
	}

	return u.q.UpdateFriendRequest(ctx, senderID, receiverID, status)
}
