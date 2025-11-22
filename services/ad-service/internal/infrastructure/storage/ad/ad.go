package ad

import (
	"context"

	"github.com/lib/pq"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"ad-service/internal/domain/entity"
	"ad-service/internal/infrastructure/storage/ad/dao"

	"ad-service/internal/pkg/transaction/wrapper"

	"ad-service/internal/pkg/aderror"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	wrapper wrapper.Database
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{wrapper: wrapper.NewDatabase(pool)}
}

func (s *Storage) CreateAd(ctx context.Context, ad *entity.Ad) error {
	var sql = `insert into public.advertisement(id,title,status,category,author_id,price,version,created_at)
			   values ($1,$2,$3,$4,$5,$6,$7,$8);`

	_, err := s.wrapper.Pool(ctx).Exec(ctx, sql,
		ad.ID(),
		ad.Title(),
		ad.Status(),
		ad.Category(),
		ad.AuthorID(),
		ad.Price(),
		ad.Version(),
		ad.CreatedAt(),
	)
	return err
}

func (s *Storage) UpdateStatus(ctx context.Context, ad *entity.Ad) error {
	var sql = `update public.advertisement set status = $1,version = $2 
               where id = $3 and version = $4;`

	commonTag, err := s.wrapper.Pool(ctx).Exec(ctx, sql,
		ad.Status(),
		ad.Version(),
		ad.ID(),
		ad.OldVersion())
	if err != nil {
		return err
	}

	if commonTag.RowsAffected() == 0 {
		return aderror.NewConflictErr("unable to update ad status due to version conflict")
	}
	return err
}

func (s *Storage) GetAd(ctx context.Context, adID uuid.UUID) (*entity.Ad, error) {
	var (
		sql = `select id,title,status,category,author_id,price,version,created_at from public.advertisement
			   where id = $1;`
		ad dao.Ad
	)

	err := pgxscan.Get(ctx, s.wrapper.Pool(ctx), &ad, sql, adID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, aderror.NewNotFoundErr("объявление не найдено", err)
		}
		return nil, err
	}

	return ad.ConvertTo(), err
}

func (s *Storage) GetAds(ctx context.Context, ids ...uuid.UUID) (entity.Ads, error) {
	var (
		sql = `select id,title,status,category,author_id,price,version,created_at from public.advertisement
			   where id =  any($1);`

		ads = make([]*dao.Ad, 0, len(ids))
	)

	err := pgxscan.Select(ctx, s.wrapper.Pool(ctx), &ads, sql, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	return lo.Map(ads, func(ad *dao.Ad, _ int) *entity.Ad {
		return ad.ConvertTo()
	}), err
}
