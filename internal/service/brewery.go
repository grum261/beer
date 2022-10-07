package service

import (
	"context"

	"github.com/grum261/beer/internal/models"
)

type BreweryRepository interface {
	CreateBrewery(ctx context.Context, p models.BreweryCreateParams) (int, error)
}

type Brewery struct {
	repo BreweryRepository
}

func NewBrewery(repo BreweryRepository) *Brewery {
	return &Brewery{
		repo: repo,
	}
}

func (b *Brewery) CreateBrewery(ctx context.Context, p models.BreweryCreateParams) (int, error) {
	return b.repo.CreateBrewery(ctx, p)
}
