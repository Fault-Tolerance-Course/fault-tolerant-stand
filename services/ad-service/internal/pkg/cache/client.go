package cache

import (
	"context"
	"fmt"
	"log/slog"

	"ad-service/internal/pkg/aderror"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

// start lifetime version
const startFrom = 1

type Client[TValue Marshaller, TValuePtr UnMarshaller] struct {
	client redis.UniversalClient
}

func NewClient[TValue Marshaller, TValuePtr UnMarshaller](client redis.UniversalClient) *Client[TValue, TValuePtr] {
	return &Client[TValue, TValuePtr]{
		client: client,
	}
}

///  ------------------  GET/SET ------------------

func (c *Client[TValue, TValuePtr]) Get(ctx context.Context, key string) (Entry[TValue, TValuePtr], error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		slog.Info(fmt.Sprintf("[Client] Cache Miss for key: < %s >", key))
		return Entry[TValue, TValuePtr]{}, errors.WithMessagef(err,
			"get entry for key %s", key)
	}

	slog.Info(fmt.Sprintf("[Client] Cache Hit for key: < %s >", key))
	return From[TValue, TValuePtr](key, value)
}

func (c *Client[TValue, TValuePtr]) Set(ctx context.Context, entry Entry[TValue, TValuePtr]) error {
	packedValue := versionPack(entry.Version, entry.marshall())
	return c.client.Set(ctx, entry.Key, packedValue, entry.TTL()).Err()
}

/// ------------------	GET/SET ------------------

/// ------------------ LOCK/UNLOCK  ------------------

func (c *Client[TValue, TValuePtr]) TryLock(ctx context.Context, entry Entry[TValue, TValuePtr]) (bool, error) {
	return c.client.SetNX(ctx, entry.Key, entry.marshall(), entry.TTL()).Result()
}

func (c *Client[TValue, TValuePtr]) Unlock(ctx context.Context, entry Entry[TValue, TValuePtr]) error {
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	keys := []string{entry.Key}
	values := []interface{}{entry.marshall()}
	cmd := script.Run(ctx, c.client, keys, values)
	return cmd.Err()
}

/// ------------------ LOCK/UNLOCK  ------------------

/// ------------------ Compare and swap V2 ------------------

// Cas модифицирует значение по ключу, увеличивая его версию
// При конкурентной ситуации выдаст ошибку
func (c *Client[TValue, TValuePtr]) Cas(ctx context.Context, entry Entry[TValue, TValuePtr]) (uint32, error) {
	var version uint32
	err := c.client.Watch(ctx, func(tx *redis.Tx) error {
		// получаем текущее значение
		packed, err := tx.Get(ctx, entry.Key).Bytes()
		// если ключа нет - мы его создадим
		if err != nil && !errors.Is(err, redis.Nil) {
			return errors.WithMessage(err, "cas: get key")
		}

		exists := !errors.Is(err, redis.Nil)
		// распаковываем версию и само значение (которое не особо то и интересно)
		version, _ = versionUnpack(packed)

		// если ключик существует и нас попросили сделать обновление (передали какую-то версию),
		// то надо понять, а можно ли это делать?
		if entry.Version != 0 && exists {
			if version != entry.Version {
				return aderror.NewConflictErr("modification conflict") // Конкурентная модификация
			}
		}

		if entry.Version == 0 && !entry.IsTombstone() && exists {
			return aderror.NewConflictErr("key unexpected exists")
		}

		// если ключа нет - его версия будет `startFrom`, в ином случае increment
		version = lo.Ternary(!exists, startFrom, version+1)

		var value []byte
		if !entry.IsTombstone() {
			value = entry.marshall()
		}

		packedValue := versionPack(version, value)

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			return pipe.Set(ctx, entry.Key, packedValue, entry.TTL()).Err()
		})
		return err
	}, entry.Key)
	if err != nil {
		if errors.Is(err, redis.TxFailedErr) {
			return 0, aderror.NewConflictErr("modification conflict") // Конкурентная модификация
		}
		return 0, err
	}
	return version, nil
}
