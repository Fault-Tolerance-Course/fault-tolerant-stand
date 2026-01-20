package ad

import (
	"context"
	"fmt"
	"log/slog"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"ad-service/internal/pkg/cache"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (a *Adapter) fromOriginOnce(ctx context.Context, adID uuid.UUID) (ad *entity.Ad, errG error) {
	// контролируем
	ad, shared, err := a.group.Do(ctx, adID.String(), func(ctx context.Context) (ad *entity.Ad, errG error) {
		// получаем данные из основного источника хранения
		ad, errG = a.origin.GetAd(ctx, adID)
		if errG != nil {
			// при ошибке выходим из функции
			slog.Error(errors.WithMessage(errG, "single flight do").Error())
			return nil, errG
		}
		return ad, nil
	})

	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("[adapter] shared key (%s): %v", adID.String(), shared))
	return ad, nil
}

func (a *Adapter) cacheOnce(ctx context.Context, entry cache.Entry[payload.Advertisement, *payload.Advertisement]) (errG error) {
	_, shared, err := a.cacheGroup.Do(ctx, entry.Key, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, a.client.Set(ctx, entry)
	})
	if err != nil {
		slog.Error(errors.WithMessage(err, "single flight do").Error())
		return err
	}

	slog.Debug(fmt.Sprintf("[adapter] shared key (%s): %v", entry.Key, shared))

	return nil
}

func (a *Adapter) refreshSkipErr(ctx context.Context, key string) {
	_, err := a.refresh(ctx, key)
	if err != nil {
		slog.Error(errors.WithMessage(err, "refresh skip err").Error())
		return
	}
	return
}

func (a *Adapter) refresh(ctx context.Context, key string) (ad *entity.Ad, err error) {
	ad, err = a.fromOriginOnce(ctx, payload.AdID(key))
	if err != nil {
		return nil, errors.WithMessage(err, "get from origin")
	}

	err = a.cacheOnce(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL, // HARD TTL
	})
	if err != nil {
		return nil, errors.WithMessage(err, "refresh ad")
	}

	return ad, nil
}
