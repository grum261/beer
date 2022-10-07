package v1

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/grum261/beer/configs"
	"github.com/grum261/beer/internal/models"
	"github.com/minio/minio-go/v7"
)

type UserService interface {
	CreateUser(ctx context.Context, p models.UserCreateParams) (int, error)
}

type UserHandler struct {
	svc         UserService
	m           *minio.Client
	maxFileSize int64
	jwtCfg      configs.JWTConfig
	v           *validator.Validate
}

func NewUserHandler(svc UserService, m *minio.Client, jwtCfg configs.JWTConfig, maxFileSize int64) *UserHandler {
	return &UserHandler{
		svc:         svc,
		m:           m,
		maxFileSize: maxFileSize,
		jwtCfg:      jwtCfg,
		v:           validator.New(),
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	err := r.ParseMultipartForm(h.maxFileSize)
	if err != nil {
		writeUnprocessable(w, err)
		return
	}

	p := models.UserCreateParams{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err = h.v.StructCtx(r.Context(), p); err != nil {
		writeUnprocessable(w, err)
		return
	}

	if bio := r.FormValue("bio"); bio != "" {
		p.Bio = &bio
	}

	f, header, err := r.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			writeUnprocessable(w, err)
			return
		}
	}
	if f != nil {
		info, err := h.m.PutObject(
			r.Context(), "users", filepath.Join("avatars", p.Username, header.Filename),
			f, header.Size, minio.PutObjectOptions{},
		)
		if err != nil {
			writeInternal(w, err)
			return
		}

		p.Avatar = &info.Key

		f.Close()
	}

	userID, err := h.svc.CreateUser(r.Context(), p)
	if err != nil {
		writeInternal(w, err)
		return
	}

	claims := models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtCfg.TokenTTL)),
		},
		Email: p.Email,
		ID:    userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtCfg.Secret))
	if err != nil {
		writeUnauthorized(w, err)
		return
	}

	writeOK(w, map[string]interface{}{
		"userID":      userID,
		"accessToken": tokenString,
		"avatar":      *p.Avatar,
	})
}
