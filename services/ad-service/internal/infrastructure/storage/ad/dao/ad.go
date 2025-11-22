package dao

import (
	"ad-service/internal/domain/entity"
	"ad-service/internal/pkg/xo"

	"github.com/shopspring/decimal"
)

type Ad struct {
	xo.Advertisement
}

func (a Ad) ConvertTo() *entity.Ad {
	return entity.With(a.ID,
		entity.WithTitle(a.Title),
		entity.WithAuthor(a.AuthorID),
		entity.WithCategory(a.Category),
		entity.WithStatus(entity.Status(a.Status)),
		entity.WithPrice(decimal.NewFromFloat(a.Price)),
		entity.WithVersion(a.Version),
		entity.WithCreatedAt(a.CreatedAt),
	)
}
