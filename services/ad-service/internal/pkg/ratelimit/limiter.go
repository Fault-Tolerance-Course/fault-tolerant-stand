package ratelimit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/samber/lo"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Limiter struct {
	state *atomic.Pointer[state]
	lock  sync.Mutex
}

func NewLimiter(config Config) *Limiter {
	stateAtomic := &atomic.Pointer[state]{}

	enriched := withDefaults(config)

	st := state{
		cfg:          enriched,
		rules:        make([]*internalState, 0, len(enriched.Rules)),
		defaultState: newInternalState(enriched.Default.MainConfig),
	}

	for _, rule := range enriched.Rules {
		st.rules = append(st.rules, newInternalState(rule))
	}

	stateAtomic.Store(&st)
	return &Limiter{
		state: stateAtomic,
	}
}

func (l *Limiter) doIfAllowed(ctx context.Context, method string, fn func() error) error {
	clientName := defaultClientName

	// fast path
	st := l.state.Load()

	if st == nil {
		return fn()
	}

	limiter, timeout := l.match(st, clientName, method)
	if limiter == nil || timeout <= 0 && limiter.Allow() {
		return fn()
	}

	if timeout > 0 {
		waitCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err := limiter.Wait(waitCtx)

		if err == nil {
			return fn()
		}

		if errors.Is(err, context.Canceled) {
			return err
		}
	}

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

	// 1. Ищем в конкретных правилах (ruleLimiters)
	for _, rule := range state.rules {
		if matcher(rule) {
			return rule.limiter, rule.cfg.Timeout
		}
	}

	// 2. Не нашли в конкретных правилах - пробуем defaultState
	defState := state.defaultState
	if defState != nil {
		if matcher(defState) {
			return defState.limiter, defState.cfg.Timeout
		}
	}

	return nil, 0
}
