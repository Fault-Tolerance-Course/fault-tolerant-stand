package ad

import (
	"context"

	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"
	"ad-service/internal/pkg/cache"

	"github.com/gofrs/uuid"
	"golang.org/x/sync/singleflight"

	singlectx "resenje.org/singleflight"
)

type Entry = cache.Entry[payload.Advertisement, *payload.Advertisement]

type Storage interface {
	GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error)

	CreateAd(ctx context.Context, ad *entity.Ad) error

	UpdateStatus(ctx context.Context, ad *entity.Ad) error
}

type Adapter struct {
	client *cache.Client[payload.Advertisement, *payload.Advertisement]
	// Источник хранения объявлений
	origin Storage
	//  Для защиты источника от лавины запросов (стандартная библиотека)
	groupStd singleflight.Group
	// Для защиты источника от лавины запросов
	group singlectx.Group[string, *entity.Ad]
	// Для защиты кеша от лавины запросов
	cacheGroup singlectx.Group[string, struct{}]
}

func NewAdapter(client *cache.Client[payload.Advertisement, *payload.Advertisement], origin Storage) *Adapter {
	return &Adapter{
		client: client,
		origin: origin,
	}
}

// GetAd возвращает объявление клиенту
func (a *Adapter) GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	return a.GetAdBetterApproach(ctx, adID)
}

func (a *Adapter) CreateAd(ctx context.Context, ad *entity.Ad) error {
	return a.CreateAdWithJitter(ctx, ad)
}

func (a *Adapter) UpdateStatus(ctx context.Context, ad *entity.Ad) error {
	return a.UpdateStatusBetterApproach(ctx, ad)
}
