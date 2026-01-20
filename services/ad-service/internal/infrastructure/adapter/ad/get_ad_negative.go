package ad

import (
	"context"
	"log/slog"
	"time"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/negative"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"ad-service/internal/pkg/aderror"
	"ad-service/internal/pkg/ctxutil"

	"github.com/gofrs/uuid"
)

// GetAdNegative возвращает объявление клиенту
func (a *Adapter) GetAdNegative(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var key = payload.AdKey(adID, config.Instance().Cache.KeyVersion)
	// получаем объявление из кеша
	entry, err := a.client.Get(ctx, key)
	// если ключ существует в кеше
	if err == nil {
		// если мы попали negative -> отдаем наружу без проверок TTL
		if negativeErr, ok := entry.Value.IsNegative(); ok {
			return nil, negativeErr
		}

		// если ключ вышел за soft ttl,
		if entry.Value.Expired(config.Instance().Cache.SoftTTL) {
			go a.refreshSkipErr(ctx, key) // то в фоне обновляем ключ (может быть отправка события!)
		}

		// Возвращаем данные, которые система считает уже stale
		return entry.Value.ConvertTo(), nil
	}

	// если произошла любая ошибка - игнорируем ее и идем в источник
	slog.Error(err.Error())

	// получаем свежие данные из источника
	ad, err := a.fromOriginOnce(ctx, adID)
	if err != nil {
		// если словили ошибку от db и она входит в список тех, что нас интересуют
		// -> кешируем негатив
		if adErr, ttl, needCache := negative.Classify(err); needCache {
			// асинхронно устанавливаем данные в кеш
			go aderror.LogWrap(a.cacheOnce(ctxutil.Detach(ctx), Entry{
				Key:        key,
				Value:      payload.ConvertFromError(adErr, time.Now()),
				Expiration: ttl, // !NEGATIVE TTL!
			}))
		}
		return nil, err
	}

	// асинхронно устанавливаем данные в кеш
	go aderror.LogWrap(a.cacheOnce(ctxutil.Detach(ctx), Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL, // HARD TTL
	}))
	return ad, nil
}
