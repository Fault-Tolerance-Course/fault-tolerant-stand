package retry

import (
	"maps"
	"sync"
	"sync/atomic"
)

type Retry struct {
	state *atomic.Pointer[stateByHandlers]
	lock  sync.Mutex
}

func NewRetry(config Config) *Retry {
	stateAtomic := &atomic.Pointer[stateByHandlers]{}

	enriched := withDefaults(config)
	state := stateByHandlers{
		cfg:            enriched,
		handlersConfig: enriched.handlersConfig(),
		handlersState:  make(map[string]*internalState),
	}

	stateAtomic.Store(&state)

	return &Retry{
		state: stateAtomic,
	}
}

func (r *Retry) getByMethod(method string) (*internalState, bool) {
	// fast path
	hedger := r.state.Load()

	state, enabled := hedger.get(method)
	if enabled && state != nil {
		return state, true
	} else if !enabled {
		return nil, false
	}

	// slow path
	r.lock.Lock()
	defer r.lock.Unlock()

	hedger = r.state.Load()

	state, enabled = hedger.get(method)
	if enabled && state != nil {
		return state, true
	} else if !enabled {
		return nil, false
	}

	// Определяем конфигурацию для метода
	var mainConfig MainConfig
	if handlerConfig, ok := hedger.handlersConfig[method]; ok {
		mainConfig = handlerConfig.MainConfig
	} else {
		mainConfig = hedger.cfg.Default.MainConfig
	}

	newHedger := *hedger
	if !*hedger.cfg.Default.SeparatePerHandlerEnabled {
		state = newInternalState(hedger.cfg.Default.MainConfig)
		newHedger.defaultState = state
	} else {
		state = newInternalState(mainConfig)

		// копируем, чтобы не было конкурентных read + write
		handlersState := maps.Clone(hedger.handlersState)
		handlersState[method] = state
		newHedger.handlersState = handlersState
	}

	r.state.Store(&newHedger)

	return state, true
}
