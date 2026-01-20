package ad

import (
	"context"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"github.com/pkg/errors"
)

// CreateAdSimple Создание объявления
func (a *Adapter) CreateAdSimple(ctx context.Context, ad *entity.Ad) error {
	err := a.origin.CreateAd(ctx, ad)
	if err != nil {
		return errors.WithMessage(err, "create ad in origin")
	}

	var key = payload.AdKey(ad.ID(), config.Instance().Cache.KeyVersion)

	err = a.client.Set(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
	})
	if err != nil {
		return errors.WithMessage(err, "set ad in cache")
	}
	return nil
}
