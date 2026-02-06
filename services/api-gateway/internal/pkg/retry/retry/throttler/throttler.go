package throttler

import (
	"sync"
)

type Throttler struct {
	max    float64 // максимальное количество токенов, которое допускаем
	thresh float64 // пороговое значение, ниже которого опускаться нельзя (иначе будет отброшен retry)
	ratio  float64 // некая доля,которую мы добавляем при успехе (токен)

	mu     sync.Mutex
	tokens float64 // текущее кол-во токенов
}

// NewThrottler создаем новый throttler
func NewThrottler(max, ratio float64) *Throttler {
	return &Throttler{
		tokens: max,     // изначально полный бак топлива
		max:    max,     // максимальный объем бака
		thresh: max / 2, // загорается предупреждение при исчерпании половины
		ratio:  ratio,
	}
}

// Throttle говорит о том, надо ли отбрасывать (throttle) retry, или нет
func (rt *Throttler) Throttle() bool {
	if rt == nil {
		return false
	}

	rt.mu.Lock()
	defer rt.mu.Unlock()

	// если перешли за минимальный порог - retry запрещен, т.к исчерпали бюджет
	if rt.tokens <= rt.thresh {
		return true
	}

	// вычитаем токен
	rt.tokens--
	// если ушли в овердрафт - закидываем в 0
	if rt.tokens < 0 {
		rt.tokens = 0
	}

	return false
}

// SuccessfulRPC засчитывает успешный rpc запрос
// Нужно для понимания того, жив ли сервис, к которому ходим. Чтобы при исчерпании квоты мы могли ее пополнить
// Ведь ранее мы лишь вычитали из нее токены
func (rt *Throttler) SuccessfulRPC() {
	if rt == nil {
		return
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	// добавляем к нашим токенам ratio (восполняем)
	rt.tokens += rt.ratio

	// ограничиваем, чтобы не переполниться хех
	if rt.tokens > rt.max {
		rt.tokens = rt.max
	}
}
