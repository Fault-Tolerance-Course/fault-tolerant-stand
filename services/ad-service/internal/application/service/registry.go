package service

import (
	"ad-service/internal/application/service/ad/create_ad"
	"ad-service/internal/application/service/ad/deactivate"
	"ad-service/internal/application/service/ad/get_active_ads"
	"ad-service/internal/application/service/ad/get_ad"
	"ad-service/internal/infrastructure/adapter"
	"ad-service/internal/infrastructure/gateway"

	"ad-service/internal/infrastructure/storage"
)

type Registry struct {
	GetAd        *get_ad.Service
	GetActiveAds *get_active_ads.Service
	Creation     *create_ad.Service
	Deactivation *deactivate.Service
}

func NewRegistry(storage *storage.Registry, gateway *gateway.Registry, adapter *adapter.Registry) *Registry {
	return &Registry{
		GetAd:        get_ad.NewService(adapter.Ad, gateway.Review),
		GetActiveAds: get_active_ads.NewService(storage.Ad),
		Creation:     create_ad.NewService(adapter.Ad),
		Deactivation: deactivate.NewService(adapter.Ad, storage.Ad),
	}
}
