package ratelimit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"ad-service/internal/pkg/grpc/clientname"

	"github.com/samber/lo"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Mode uint8

const (
	ModeClient = iota
	ModeServer
)

type Limiter struct {
	state *atomic.Pointer[state]
	lock  sync.Mutex
}

func NewLimiter(config Config) *Limiter {
	stateAtomic := &atomic.Pointer[state]{}

	st := state{
		cfg:   config,
		rules: make([]*internalState, 0, len(config.Rules)),
	}

	for _, rule := range config.Rules {
		st.rules = append(st.rules, newInternalState(rule))
	}

	stateAtomic.Store(&st)
	return &Limiter{
		state: stateAtomic,
	}
}

func (l *Limiter) doIfAllowed(ctx context.Context, method string, fn func() error) error {
	// fast path
	st := l.state.Load()
	if st == nil {
		return fn()
	}

	clientName := clientname.FromContext(ctx)

	limiter, timeout := l.match(st, clientName, method)
	if limiter == nil || timeout <= 0 && limiter.Allow() {
		return fn()
	}

	// ВСЕГДА ставим Header сразу — состояние на входе
	setHeader(ctx, limiter)

	if timeout > 0 {
		waitCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err := limiter.Wait(waitCtx)

		if err == nil {
			// ждали, но прошли
			setTrailers(ctx, 0)
			return fn()
		}

		if errors.Is(err, context.Canceled) {
			return err
		}
	}

	// высчитываем время до следующего токена
	setTrailers(ctx, retryAfter(limiter))
	return status.Errorf(
		codes.ResourceExhausted,
		"rate limit exceeded for %#q client, current max rate per sec.: %d", clientName, int(limiter.Limit()),
	)
}

func (l *Limiter) match(state *state, client, handler string) (*rate.Limiter, time.Duration) {
	var matcher = func(internal *internalState) bool {
		var (
			clientMatch  = len(internal.cfg.Clients) == 0 || lo.Contains(internal.cfg.Clients, client)
			handlerMatch = len(internal.cfg.Handlers) == 0 || lo.Contains(internal.cfg.Handlers, handler)
		)

		return clientMatch && handlerMatch
	}

	for _, rule := range state.rules {
		if matcher(rule) {
			return rule.limiter, rule.cfg.Timeout
		}
	}
	return nil, 0
}

func retryAfter(l *rate.Limiter) time.Duration {
	now := time.Now()

	r := l.ReserveN(now, 1)
	if !r.OK() {
		return rate.InfDuration
	}

	delay := r.DelayFrom(now)
	r.CancelAt(now) // ОБЯЗАТЕЛЬНО делаем Cancel, чтобы токен вернулся и потерялся

	return delay
}
