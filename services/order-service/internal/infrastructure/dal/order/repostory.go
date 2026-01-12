package order

import (
	"context"
	"errors"

	"order-service/internal/domain/entity"
	"order-service/internal/domain/value_object"
	"order-service/internal/infrastructure/dal/order/dao"
	"order-service/internal/pkg/ordererror"
	"order-service/internal/pkg/transaction/wrapper"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	wrapper wrapper.Database
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{wrapper: wrapper.NewDatabase(pool)}
}

func (r *Repository) GetOrdersVersion(ctx context.Context, clientID uuid.UUID) (value_object.Version, error) {
	var sql = `select version from orders_state where client_id =$1;`

	var version int64
	err := pgxscan.Get(ctx, r.wrapper.Pool(ctx), &version, sql, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// скажем, я решил создавать тут запись при создании 1-го заказа
			return value_object.Version{Old: 0}, nil
		}
	}

	return value_object.Version{Old: version}, nil
}

func (r *Repository) UpdateOrdersVersion(ctx context.Context, clientID uuid.UUID, version value_object.Version) error {
	var sql = `insert into orders_state (client_id, version)
			   values ($1, $3)
			   on conflict (client_id) do update
			   set version = excluded.version
			   where orders_state.version = $2;`

	commonTag, err := r.wrapper.Pool(ctx).Exec(ctx, sql, clientID, version.Previous(), version.Current())
	if err != nil {
		return err
	}

	if commonTag.RowsAffected() == 0 {
		return ordererror.NewConflictErr("unable to update orders state due to version conflict")
	}

	return nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *entity.Order, key string) error {
	var sql = `insert into public.order (
                    order_id,
                    idempotency_key,
                    status,
                    client_id,
                    price,
                    ad_id,
                    category,
                    created_at
                    )
				values ($1,$2,$3,$4,$5,$6,$7,$8)
				on conflict (idempotency_key) do nothing;`
	commonTag, err := r.wrapper.Pool(ctx).Exec(ctx, sql,
		order.ID(),
		key,
		order.Status(),
		order.ClientID(),
		order.Price(),
		order.AdID(),
		order.Category(),
		order.CreatedAt(),
	)

	if err != nil {
		return err
	}

	// конфликт по ключу идемпотентности
	if commonTag.RowsAffected() == 0 {
		return ordererror.NewConflictErr("unable to create order due to idempotency conflict")
	}

	return nil
}

func (r *Repository) GetActiveOrders(ctx context.Context, clientID uuid.UUID) (entity.OrderList, error) {
	var (
		sql = `
		select
			o.client_id,
			coalesce(s.version, 0) as orders_version,
			o.order_id,
			o.status,
			o.price,
			o.ad_id,
			o.category,
			o.created_at
		from public."order" as o
		left join orders_state as s on s.client_id = o.client_id
		where o.client_id = $1;`

		orders = make(dao.Orders, 0)
	)

	err := pgxscan.Select(ctx, r.wrapper.Pool(ctx), &orders, sql, clientID)
	if err != nil {
		return entity.OrderList{}, err
	}

	return orders.Convert(), nil
}
