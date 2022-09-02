package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

func NewUserHandler(svc UserService, jwtCfg configs.JWTConfig, maxFileSize int64) *UserHandler {
	return &UserHandler{
		svc:         svc,
		m:           &minio.Client{},
		maxFileSize: maxFileSize,
		jwtCfg:      jwtCfg,
		v:           validator.New(),
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r http.Request, params map[string]string) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(h.maxFileSize); err != nil {
		_ = enc.Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := models.UserCreateParams{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err := h.v.StructCtx(r.Context(), p); err != nil {
		_ = enc.Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bio := r.FormValue("bio"); bio != "" {
		p.Bio = &bio
	}

	f, header, err := r.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			_ = enc.Encode(map[string]string{
				"error": err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if f != nil {
		info, err := h.m.PutObject(
			r.Context(), "users", "avatars/"+p.Username,
			f, header.Size, minio.PutObjectOptions{},
		)
		if err != nil {
			_ = enc.Encode(map[string]string{
				"error": err.Error(),
			})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p.Avatar = &info.Location

		f.Close()
	}

	userID, err := h.svc.CreateUser(r.Context(), p)
	if err != nil {
		_ = enc.Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
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
		_ = enc.Encode(map[string]string{
			"error": err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = enc.Encode(map[string]interface{}{
		"userID":      userID,
		"accessToken": tokenString,
		"avatar":      *p.Avatar,
	})
	w.WriteHeader(http.StatusCreated)
}
