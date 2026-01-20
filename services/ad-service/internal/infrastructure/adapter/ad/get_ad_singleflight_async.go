package ad

import (
	"context"
	"log/slog"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"ad-service/internal/pkg/aderror"
	"ad-service/internal/pkg/ctxutil"

	"github.com/gofrs/uuid"
)

// GetAdSingleFlightAsync возвращает объявление клиенту
func (a *Adapter) GetAdSingleFlightAsync(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var key = payload.AdKey(adID, config.Instance().Cache.KeyVersion)
	// получаем объявление из кеша
	entry, err := a.client.Get(ctx, key)
	// если ключ существует в кеше
	if err == nil {
		return entry.Value.ConvertTo(), nil
	}

	// если произошла любая ошибка - игнорируем ее и идем в источник
	slog.Error(err.Error())

	// получаем свежие данные из источника
	ad, err := a.fromOriginOnce(ctx, adID)
	if err != nil {
		return nil, err
	}

	// асинхронно устанавливаем данные в кеш
	go aderror.LogWrap(a.cacheOnce(ctxutil.Detach(ctx), Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
	}))

	return ad, nil
}
