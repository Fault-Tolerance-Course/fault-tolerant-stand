package payload

import (
	"time"

	"ad-service/internal/domain/entity"
	"ad-service/internal/pkg/aderror"
	"ad-service/internal/pkg/cache"

	"ad-service/internal/infrastructure/adapter/ad/negative"

	"github.com/gofrs/uuid"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var Tombstone = cache.NewTombstone[Advertisement, *Advertisement]

//go:generate easyjson -all payload.go

//easyjson:json
type Advertisement struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	AuthorID  string    `json:"author_id"`
	Price     string    `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CachedAt  time.Time `json:"cached_at"`
	// В случае неуспеха можно кешировать негативный ответ
	Negative negative.Payload `json:"negative,omitempty"`
}

func (a *Advertisement) Expired(ttl time.Duration) bool {
	return time.Since(a.CachedAt) >= ttl
}

func (a *Advertisement) IsNegative() (err error, negative bool) {
	if lo.IsEmpty(a.Negative) {
		return nil, false
	}
	return aderror.New(
		aderror.WithCode(a.Negative.Code),
		aderror.WithClassification(a.Negative.Classification),
		aderror.WithMessage(a.Negative.Message),
		aderror.WithBase(a.Negative.BaseMessage),
	), true
}

func ConvertFromError(err *aderror.Error, cachedAt time.Time) Advertisement {
	return Advertisement{
		Negative: negative.Payload{
			Code:           int(err.Code()),
			Classification: err.Classification(),
			Message:        err.Message(),
			BaseMessage:    err.BaseMessage(),
		},
		CachedAt: cachedAt,
	}
}

func ConvertFrom(ad *entity.Ad) Advertisement {
	return Advertisement{
		ID:        ad.ID().String(),
		Title:     ad.Title(),
		Category:  ad.Category(),
		AuthorID:  ad.AuthorID().String(),
		Price:     ad.Price().String(),
		Status:    string(ad.Status()),
		CreatedAt: ad.CreatedAt(),
		CachedAt:  time.Now(), // для небольшой сигнатуры оставлю тут
	}
}

func (a *Advertisement) ConvertTo() *entity.Ad {
	return entity.With(uuid.FromStringOrNil(a.ID),
		entity.WithTitle(a.Title),
		entity.WithAuthor(uuid.FromStringOrNil(a.ID)),
		entity.WithCategory(a.Category),
		entity.WithStatus(entity.Status(a.Status)),
		entity.WithPrice(decimal.RequireFromString(a.Price)),
		entity.WithCreatedAt(a.CreatedAt),
	)
}
