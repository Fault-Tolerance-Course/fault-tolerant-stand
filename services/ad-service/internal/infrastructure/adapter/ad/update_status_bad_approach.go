package ad

import (
	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"context"

	"github.com/pkg/errors"
)

func (a *Adapter) UpdateStatusBadApproach(ctx context.Context, ad *entity.Ad) error {
	err := a.origin.UpdateStatus(ctx, ad)
	if err != nil {
		return errors.WithMessage(err, "create ad in origin")
	}

	var key = payload.AdKey(ad.ID(), config.Instance().Cache.KeyVersion)

	// А если будет неуспех?
	// Плохо это
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
