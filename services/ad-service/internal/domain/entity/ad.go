package entity

import (
	"time"

	"ad-service/internal/domain/value_object"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"

	"github.com/shopspring/decimal"
)

type Status string

const (
	// StatusActive статус активное
	StatusActive Status = "active"
	// StatusInactive статус неактивное
	StatusInactive Status = "inactive"
)

type version struct {
	Old int
	new *int
}

func (v *version) inc() {
	v.new = lo.ToPtr(v.Old + 1)
}

// Ad представление доменной модели агрегата
type Ad struct {
	id        uuid.UUID
	title     string
	category  string
	authorID  uuid.UUID
	price     decimal.Decimal
	status    Status
	createdAt time.Time

	review *value_object.Review

	version version
}

func With(id uuid.UUID, opts ...Option) *Ad {
	ad := &Ad{id: id}

	for _, opt := range opts {
		opt(ad)
	}
	return ad
}
func NewAd(title, category string, authorID uuid.UUID, price decimal.Decimal) *Ad {
	id, _ := uuid.NewV7()
	return &Ad{
		id:        id,
		title:     title,
		category:  category,
		authorID:  authorID,
		price:     price,
		status:    StatusActive,
		createdAt: time.Now(),
	}
}

func (a *Ad) ID() uuid.UUID {
	return a.id
}

func (a *Ad) Title() string {
	return a.title
}

func (a *Ad) Category() string {
	return a.category
}

func (a *Ad) AuthorID() uuid.UUID {
	return a.authorID
}

func (a *Ad) Price() decimal.Decimal {
	return a.price
}

func (a *Ad) Status() Status {
	return a.status
}

func (a *Ad) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Ad) Deactivate() {
	a.status = StatusInactive
	a.version.inc()
}

func (a *Ad) Version() int {
	if a.version.new != nil {
		return *a.version.new
	}
	return a.version.Old
}

func (a *Ad) OldVersion() int {
	return a.version.Old
}

func (a *Ad) SetReview(review decimal.Decimal) {
	a.review = lo.ToPtr(value_object.Review{
		Avg: review,
	})
}

type Ads []*Ad

func (as Ads) Active() []*Ad {
	return lo.Filter(as, func(a *Ad, _ int) bool {
		return a.Status() == StatusActive
	})
}
