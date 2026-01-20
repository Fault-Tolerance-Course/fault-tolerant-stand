package ad

import (
	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// GetAdSimple возвращает объявление клиенту
func (a *Adapter) GetAdSimple(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var key = payload.AdKey(adID, config.Instance().Cache.KeyVersion)
	// получаем объявление из кеша
	entry, err := a.client.Get(ctx, key)
	// если ключ существует в кеше
	if err == nil {
		return entry.Value.ConvertTo(), nil
	}

	// если произошла любая ошибка - игнорируем ее и идем в источник
	slog.Error(err.Error())

	// получаем данные из основного источника хранения
	ad, err := a.origin.GetAd(ctx, adID)
	if err != nil {
		// при ошибке выходим из функции
		slog.Error(errors.WithMessage(err, "get ad from origin").Error())
		return nil, err
	}

	err = a.client.Set(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
	})
	if err != nil {
		slog.Error(errors.WithMessage(err, "set ad to cache").Error())
		return nil, err
	}

	return ad, err
}
