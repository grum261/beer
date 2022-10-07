package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/grum261/beer/internal/models"
	"github.com/pkg/errors"
)

type UserQueries struct {
	*Queries
}

func (q *UserQueries) InsertUser(ctx context.Context, p models.UserCreateParams) (int, error) {
	query, args, err := q.stmtBuilder.
		Insert("users").
		Columns(
			"username",
			"email",
			"password_hash",
			"bio",
			"avatar").
		Values(
			p.Username, p.Email,
			p.Password, p.Bio,
			p.Avatar,
		).ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "postgres.InsertUser")
	}

	var userID int

	if err := q.db.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
		return userID, errors.Wrap(err, "postgres.InsertUser")
	}

	return userID, nil
}

func (q *UserQueries) SelectPassword(ctx context.Context, arg, field string) (int, string, error) {
	query, args, err := q.stmtBuilder.
		Select(
			"id",
			"password_hash",
		).
		From("users").
		Where(sq.Eq{field: arg}).
		ToSql()
	if err != nil {
		return 0, "", errors.Wrap(err, "postgres.SelectPassword")
	}

	var (
		hashedPass string
		userID     int
	)

	if err := q.db.QueryRow(ctx, query, args...).Scan(&userID, &hashedPass); err != nil {
		return userID, hashedPass, errors.Wrap(err, "postgres.SelectPassword")
	}

	return userID, hashedPass, nil
}

func (q *UserQueries) InsertFriendRequest(ctx context.Context, senderID, receiverID int) error {
	query, args, err := q.stmtBuilder.
		Insert("users_friends_requests").
		Columns("user_sender_id", "user_receiver_id", "request_status").
		Values(senderID, receiverID, models.StatusSent).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres.InsertFriendRequest")
	}

	if _, err := q.db.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(err, "postgres.InsertFriendRequest")
	}

	return nil
}

func (q *UserQueries) UpdateFriendRequest(ctx context.Context, senderID, receiverID int, status models.RequestStatus) error {
	query, args, err := q.stmtBuilder.
		Update("users_friends_requests").
		SetMap(map[string]interface{}{
			"request_status": status,
			"updated_at":     time.Now(),
		}).
		Where(sq.And{
			sq.Eq{"user_sender_id": senderID},
			sq.Eq{"user_receiver_id": receiverID},
		}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres.UpdateFriendRequest")
	}

	res, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "postgres.UpdateFriendRequest")
	}

	if res.RowsAffected() == 0 {
		return errors.New("postgres.UpdateFriendRequest: failed to find friend request")
	}

	return nil
}

func (q *UserQueries) SelectUserFriends(ctx context.Context, userID int) ([]models.UserFriend, error) {
	query, args, err := q.stmtBuilder.Select(
		"u.id",
		"u.username",
		"u.avatar",
		"ur.updated_at",
	).
		From("users u").
		InnerJoin("users_friends_requests ur ON u.id = ur.receiver_id").
		Where(sq.And{
			sq.Eq{"u.id": userID},
			sq.Eq{"ur.request_status": models.StatusAccepted},
		}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "postgres.SelectUserFriends")
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "postgres.SelectUserFriends")
	}
	defer rows.Close()

	var usersFriends []models.UserFriend

	for rows.Next() {
		var friend models.UserFriend

		if err := rows.Scan(
			&friend.ID,
			&friend.Username,
			&friend.Avatar,
			&friend.FriendsSince,
		); err != nil {
			return nil, errors.Wrap(err, "postgres.SelectUserFriends")
		}

		usersFriends = append(usersFriends, friend)
	}

	return usersFriends, nil
}

func (q *UserQueries) SelectActiveFriendsRequests(ctx context.Context, userID int) ([]models.ActiveFriendRequest, error) {
	q.stmtBuilder.Select(
		"u.id",
		"u.username",
		"u.avatar",
		"ur.sent_at",
	)

	return nil, nil
}
