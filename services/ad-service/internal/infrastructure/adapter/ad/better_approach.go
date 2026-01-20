package ad

import (
	"context"
	"fmt"
	"log/slog"

	"ad-service/config"
	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/adapter/ad/payload"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

/*
	Запрос 1 (пишет статус)                            Запрос 2 (пишет категорию)
	--------------------------------------------------------------------------
	1. v2 = cas(tombstone) v1 -> v2
													   2. v3 = cas(tombstone) v2 -> v3
	--------------------------------------------------------------------------
	3. update_status(ad)
													   4. update_category(ad)
	--------------------------------------------------------------------------
	5. cas(new_value,v2) -> fail (v3)
       												   6. cas(new_value,v3) v3->v4 = success
*/

// Почему при tombstone нельзя делать lazy loading?
/*
	Запрос 1 (пишет статус)                            Запрос 2 (читает объявление)
	--------------------------------------------------------------------------
	1. v2 = cas(tombstone) v1 -> v2
													   2. get(key) -> v2:tombstone
	--------------------------------------------------------------------------
													   3. read_db() -> old_data
	4. update_status(new_data)
	--------------------------------------------------------------------------
			                                           5. cas(old_data,v2) v2->v3
    6. cas(new_data,v2) -> fail (v3)
*/

/*
	Запрос 1 (пишет статус)                            Запрос 2 (читает объявление)
	--------------------------------------------------------------------------
	 												   1. get(key) -> redis.Nil
	2. cas(tombstone) nil -> v1
	--------------------------------------------------------------------------
													   3. read_db() -> old_data
	4. update_status(new_data)
	--------------------------------------------------------------------------
													   5. cas(old_data,0) -> fail
    6. cas(new_data,v1) v1->v2
*/

// Какой у нас алгоритм действий?
/*
	1. DEL key
 	2. UPDATE DATABASE
	3. SET key value
*/

func (a *Adapter) UpdateStatusBetterApproach(ctx context.Context, ad *entity.Ad) error {
	// Уже знаем, что версия сущности entity.Ad увеличена
	var key = payload.AdKey(ad.ID(), config.Instance().Cache.KeyVersion)
	// 1. Удаляем ключ, чтобы в ситуации, когда на последнем шаге
	// не записали в кеш - данные не остались stale
	version, err := a.client.Cas(ctx, payload.Tombstone(key, config.Instance().Cache.TTL))
	if err != nil {
		return errors.Wrap(err, "cas invalidate cache")
	}

	// 2. Выполняем обновление на актуальное состояние в базе благодаря optimistic lock
	err = a.origin.UpdateStatus(ctx, ad)
	if err != nil {
		return errors.WithMessage(err, "update ad status in origin")
	}

	// 3. Выполняем установку актуального значения в кеш благодаря версии, полученной ранее (токен)
	version, err = a.client.Cas(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
		Version:    version,
	})
	if err != nil {
		return errors.WithMessage(err, "cas ad in cache")
	}

	slog.Info(fmt.Sprintf("updated ad versions: [origin]: %d, [cache]: %d",
		ad.Version(), version))
	return nil
}

func (a *Adapter) GetAdBetterApproach(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var key = payload.AdKey(adID, config.Instance().Cache.KeyVersion)
	// получаем объявление из кеша
	entry, err := a.client.Get(ctx, key)
	// если ключ существует в кеше и он не tombstone
	if err == nil && !entry.IsTombstone() {
		return entry.Value.ConvertTo(), nil
	}

	// если произошла любая ошибка - игнорируем ее и идем в источник
	if err != nil {
		slog.Error(err.Error())
	}

	// На данном этапе ключ либо tombstone, либо его нет/ошибка

	// получаем свежие данные из источника
	ad, err := a.fromOriginOnce(ctx, adID)
	if err != nil {
		return nil, err
	}

	// если ключ Tombstone - это индикатор того,
	// что сейчас идет запись и мы не будем этого делать здесь
	if entry.IsTombstone() {
		return ad, nil
	}

	version, err := a.client.Cas(ctx, Entry{
		Key:        key,
		Value:      payload.ConvertFrom(ad),
		Expiration: config.Instance().Cache.TTL,
		Version:    entry.Version, // Берем токен после получения (если его нет - он создастся)
	})

	slog.Info(fmt.Sprintf("updated ad versions: [cache]: %d", version))
	return ad, nil
}
