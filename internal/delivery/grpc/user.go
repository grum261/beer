package grpc

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/auth"
	"github.com/grum261/beer/configs"
	"github.com/grum261/beer/internal/models"
	"github.com/grum261/beer/proto/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	AuthUser(ctx context.Context, p models.UserAuthParams) (int, bool, error)
	SendFriendRequest(ctx context.Context, senderID, receiverID int) error
	AcceptFriendRequest(ctx context.Context, senderID, receiverID int) error
	DeclineFriendRequest(ctx context.Context, senderID, receiverID int) error
}

type UserDelivery struct {
	svc    UserService
	jwtCfg configs.JWTConfig
	userpb.UserDeliveryServiceServer
}

func NewUserDelivery(svc UserService, jwtConfig configs.JWTConfig) *UserDelivery {
	return &UserDelivery{
		svc:    svc,
		jwtCfg: jwtConfig,
	}
}

func (u *UserDelivery) AuthUserHandler(
	ctx context.Context,
	in *userpb.AuthUserRequest,
) (*userpb.AuthUserResponse, error) {
	userID, ok, err := u.svc.AuthUser(ctx, models.UserAuthParams{
		Username: in.User.Username,
		Email:    in.User.Email,
		Password: in.User.Password,
	})
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"invalid password for user: %v %v", in.User.Username, in.User.Email,
		)
	}

	claims := models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.jwtCfg.TokenTTL)),
		},
		Email: in.User.Email,
		ID:    userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtCfg.Secret))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.AuthUserResponse{
		AccessToken: tokenString,
	}, nil
}

func (u *UserDelivery) SendFriendRequestHandler(
	ctx context.Context,
	in *userpb.FriendRequest,
) (*userpb.FriendResponse, error) {
	sender, ok := ctx.Value(tokenKey{}).(*models.UserClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid claims")
	}

	if err := u.svc.SendFriendRequest(ctx, sender.ID, int(in.ReceiverId)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.FriendResponse{
		IsSended: true,
	}, nil
}

func (u *UserDelivery) UpdateFriendRequestHandler(
	ctx context.Context,
	in *userpb.UpdateFriendRequest,
) (*emptypb.Empty, error) {
	sender, ok := ctx.Value(tokenKey{}).(*models.UserClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid claims")
	}

	var err error

	switch {
	case in.Status == 0:
		err = u.svc.DeclineFriendRequest(ctx, sender.ID, int(in.ReceiverId))
	case in.Status == 1:
		err = u.svc.AcceptFriendRequest(ctx, sender.ID, int(in.ReceiverId))
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid request status")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (u *UserDelivery) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	switch {
	case strings.Contains(fullMethodName, "CreateUserHandler"),
		strings.Contains(fullMethodName, "AuthUserHandler"):
		return ctx, nil
	}

	authFunc := JWTAuth(u.jwtCfg.Secret)

	return authFunc(ctx)
}

type tokenKey struct{}

func JWTAuth(secret string) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		jwtToken, err := jwt.ParseWithClaims(
			token, &models.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, errors.New("unexpected token signature")
				}

				return []byte(secret), nil
			},
		)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		claims, ok := jwtToken.Claims.(*models.UserClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		newCtx := context.WithValue(ctx, tokenKey{}, claims)

		return newCtx, nil
	}
}
