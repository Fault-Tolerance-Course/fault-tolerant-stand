package ad

import (
	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// GetAdSingleFlightSTD возвращает объявление клиенту
func (a *Adapter) GetAdSingleFlightSTD(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var key = payload.AdKey(adID, config.Instance().Cache.KeyVersion)
	// получаем объявление из кеша
	entry, err := a.client.Get(ctx, key)
	// если ключ существует в кеше
	if err == nil {
		return entry.Value.ConvertTo(), nil
	}

	// если произошла любая ошибка - игнорируем ее и идем в источник
	slog.Error(err.Error())

	var (
		ad     *entity.Ad
		shared bool
	)
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// контролируем
	resCh := a.groupStd.DoChan(adID.String(), func() (interface{}, error) {
		// получаем данные из основного источника хранения (контекст передан, все в порядке)
		ad, err = a.origin.GetAd(ctx, adID)
		if err != nil {
			// при ошибке выходим из функции
			slog.Error(errors.WithMessage(err, "single flight do").Error())
			return nil, err
		}

		return ad, err
	})

	select {
	case <-ctx.Done():
		a.groupStd.Forget(adID.String())
		return nil, ctx.Err()
	case res := <-resCh:
		// Забываем при ошибках (кроме отмены контекста)
		if res.Err != nil {
			if !errors.Is(res.Err, context.Canceled) {
				a.groupStd.Forget(key)
			}
			return nil, res.Err
		}
		ad, shared = res.Val.(*entity.Ad), res.Shared
	}

	slog.Debug(fmt.Sprintf("[adapter] shared key (%s): %v", adID.String(), shared))

	// обновляем в cache
	err = a.client.Set(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
	})
	if err != nil {
		// при ошибке
		slog.Error(errors.WithMessage(err, "set ad to cache").Error())
	}
	return ad, nil
}
