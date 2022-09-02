package models

import "github.com/golang-jwt/jwt/v4"

type UserCreateParams struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
	Bio      *string
	Avatar   *string
}

type UserAuthParams struct {
	Username string `validate:"omitempty"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    int
}
