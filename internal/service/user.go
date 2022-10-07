package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/grum261/beer/configs"
	"github.com/grum261/beer/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

type UserRepository interface {
	CreateUser(ctx context.Context, p models.UserCreateParams) (int, error)
	GetUserPasswordByName(ctx context.Context, username string) (int, string, error)
	GetUserPasswordByEmail(ctx context.Context, email string) (int, string, error)
	UpdateFriendStatus(ctx context.Context, senderID, receiverID int, status models.RequestStatus) error
}

type User struct {
	repo UserRepository
	v    *validator.Validate
	cfg  configs.Argon2Config
}

func NewUserService(repo UserRepository, cfg configs.Argon2Config) *User {
	return &User{
		repo: repo,
		v:    validator.New(),
		cfg:  cfg,
	}
}

func (u *User) CreateUser(ctx context.Context, p models.UserCreateParams) (int, error) {
	if err := u.v.StructCtx(ctx, p); err != nil {
		return 0, errors.Wrap(err, "service.CreateUser: validation failed")
	}

	hashedPass, err := u.generateHash(p.Password)
	if err != nil {
		return 0, errors.Wrap(err, "service.CreateUser: hash generation failed")
	}

	p.Password = hashedPass

	userID, err := u.repo.CreateUser(ctx, p)
	if err != nil {
		return userID, errors.Wrap(err, "service.CreateUser: user creation failed")
	}

	return userID, nil
}

func (u *User) AuthUser(ctx context.Context, p models.UserAuthParams) (int, bool, error) {
	err := u.v.StructCtx(ctx, p)
	if err != nil {
		return 0, false, errors.Wrap(err, "service.AuthUser: validation failed")
	}

	var (
		hashedPass string
		userID     int
	)

	switch {
	case p.Username != "":
		userID, hashedPass, err = u.repo.GetUserPasswordByName(ctx, p.Username)
	case p.Email != "":
		userID, hashedPass, err = u.repo.GetUserPasswordByEmail(ctx, p.Email)
	}
	if err != nil {
		return 0, false, errors.Wrap(err, "service.AuthUser: getting user failed")
	}

	ok, err := comparePassAndHash(p.Password, hashedPass)
	if err != nil {
		return userID, ok, errors.Wrap(err, "service.AuthUser: invalid password")
	}

	return userID, ok, nil
}

func (u *User) SendFriendRequest(ctx context.Context, senderID, receiverID int) error {
	if err := u.repo.UpdateFriendStatus(
		ctx, senderID, receiverID, models.StatusSent,
	); err != nil {
		return errors.Wrap(err, "service.SendFriendRequest: friend request status update failed")
	}

	return nil
}

func (u *User) AcceptFriendRequest(ctx context.Context, senderID, receiverID int) error {
	if err := u.repo.UpdateFriendStatus(
		ctx, senderID, receiverID, models.StatusAccepted,
	); err != nil {
		return errors.Wrap(err, "service.AcceptFriendRequest: friend request status update failed")
	}

	return nil
}

func (u *User) DeclineFriendRequest(ctx context.Context, senderID, receiverID int) error {
	if err := u.repo.UpdateFriendStatus(
		ctx, senderID, receiverID, models.StatusDeclined,
	); err != nil {
		return errors.Wrap(err, "service.DeclineFriendRequest: friend request status update failed")
	}

	return nil
}

func (u *User) generateHash(pass string) (string, error) {
	salt := make([]byte, u.cfg.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(pass), salt, u.cfg.Iterations,
		u.cfg.Memory, u.cfg.Threads, u.cfg.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, u.cfg.Memory, u.cfg.Iterations,
		u.cfg.Threads, b64Salt, b64Hash,
	)

	return encodedHash, nil
}

func comparePassAndHash(pass, encodedPass string) (bool, error) {
	p, salt, hash, err := decodeHash(encodedPass)
	if err != nil {
		return false, err
	}

	other := argon2.IDKey([]byte(pass), salt, p.Iterations, p.Memory, p.Threads, p.KeyLength)

	return subtle.ConstantTimeCompare(hash, other) == 1, nil
}

func decodeHash(encodedPass string) (configs.Argon2Config, []byte, []byte, error) {
	vals := strings.Split(encodedPass, "$")
	if len(vals) != 6 {
		return configs.Argon2Config{}, nil, nil, errors.New("the encoded hash is not in the correct format")
	}

	var ver int
	_, err := fmt.Sscanf(vals[2], "v=%d", &ver)
	if err != nil {
		return configs.Argon2Config{}, nil, nil, err
	}

	if ver != argon2.Version {
		return configs.Argon2Config{}, nil, nil, errors.New("incompatible version of argon2")
	}

	p := configs.Argon2Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Threads)
	if err != nil {
		return configs.Argon2Config{}, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return configs.Argon2Config{}, nil, nil, err
	}

	p.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return configs.Argon2Config{}, nil, nil, err
	}

	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
