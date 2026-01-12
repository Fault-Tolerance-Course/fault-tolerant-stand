package hedge

import (
	"sync"
	"time"
)

type DurationBuf struct {
	mu      sync.RWMutex
	data    []time.Duration
	current int
}

func NewDurationBuf(size int) *DurationBuf {
	return &DurationBuf{
		data: make([]time.Duration, size),
	}
}

func (s *DurationBuf) Append(v time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[s.current] = v
	s.current++

	if s.current >= len(s.data) {
		s.current = 0
	}
}

func (s *DurationBuf) Copy() []time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stat := make([]time.Duration, len(s.data))

	copy(stat, s.data)

	return stat
}
