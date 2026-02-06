package retry

import (
	"profile-service/internal/pkg/retry/throttler"
)

type stateByHandlers struct {
	cfg            Config
	handlersConfig map[string]HandlerConfig
	handlersState  map[string]*internalState
	defaultState   *internalState
}

type internalState struct {
	cfg MainConfig

	throttler *throttler.Throttler
}

func newInternalState(config MainConfig) *internalState {
	var th *throttler.Throttler
	if config.Throttler.Enabled {
		th = throttler.NewThrottler(config.Throttler.MaxTokens, config.Throttler.TokenRatio)
	}

	return &internalState{
		cfg:       config,
		throttler: th,
	}
}

func (s stateByHandlers) get(method string) (_ *internalState, enabled bool) {
	st, ok := s.handlersState[method]
	if ok {
		return st, true
	}

	isEnabled := s.isHandlerEnabled(method)
	if !isEnabled {
		return nil, false
	}

	if !*s.cfg.Default.SeparatePerHandlerEnabled && s.defaultState != nil {
		return s.defaultState, true
	}

	return nil, true
}

func (s stateByHandlers) isHandlerEnabled(method string) bool {
	handlerConfig, ok := s.handlersConfig[method]
	if ok {
		return handlerConfig.Enabled
	}

	return s.cfg.Default.Enabled
}
