package ad

import (
	"context"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"
	"ad-service/internal/pkg/jitter"

	"github.com/pkg/errors"
)

func (a *Adapter) CreateAdWithJitter(ctx context.Context, ad *entity.Ad) error {
	err := a.origin.CreateAd(ctx, ad)
	if err != nil {
		return errors.WithMessage(err, "create ad in origin")
	}

	var key = payload.AdKey(ad.ID(), config.Instance().Cache.KeyVersion)

	// вычисляем с jitter +- 20% от основного TTL
	ttl := jitter.UniformPositiveJitter(
		config.Instance().Cache.TTL, // Основной TTL
		config.Instance().Cache.JitterPercentage,
	)

	err = a.client.Set(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: ttl,
	})
	if err != nil {
		return errors.WithMessage(err, "set ad in cache")
	}
	return nil
}
