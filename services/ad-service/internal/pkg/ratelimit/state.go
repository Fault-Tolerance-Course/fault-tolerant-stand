package ratelimit

import "golang.org/x/time/rate"

const (
	defaultBurstRatio = 1.0
)

type state struct {
	cfg Config

	rules []*internalState
}

type internalState struct {
	// Фактически правило
	cfg MainConfig
	// Ограничитель
	limiter *rate.Limiter
}

func newInternalState(config MainConfig) *internalState {
	// инициализируем Burst
	burst := int(config.Limit * defaultBurstRatio)
	if config.Burst > 0 {
		burst = config.Burst
	}

	// создаем ограничитель
	limiter := rate.NewLimiter(rate.Limit(config.Limit), burst)

	return &internalState{
		cfg:     config,
		limiter: limiter,
	}
}
