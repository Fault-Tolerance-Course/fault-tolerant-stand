package outbox

type Option func(outbox *Outbox)

func WithMetrics() Option {
	return func(_ *Outbox) {
		registerMetrics()
	}
}
