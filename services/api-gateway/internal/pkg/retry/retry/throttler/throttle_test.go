package throttler

import "testing"

func TestThrottler(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Throttler
		actions func(*Throttler) (throttle bool)
		want    bool
	}{
		{
			name: "initial → no throttle",
			setup: func() *Throttler {
				return NewThrottler(10, 1)
			},
			actions: func(th *Throttler) bool {
				return th.Throttle()
			},
			want: false,
		},
		{
			name: "exhaust tokens → throttle",
			setup: func() *Throttler {
				return NewThrottler(2, 1)
			},
			actions: func(th *Throttler) bool {
				th.Throttle()        // 1
				th.Throttle()        // 0
				return th.Throttle() // <= thresh → throttle
			},
			want: true,
		},
		{
			name: "SuccessfulRPC refills tokens → no throttle",
			setup: func() *Throttler {
				return NewThrottler(3, 2)
			},
			actions: func(th *Throttler) bool {
				th.Throttle()        // 2
				th.Throttle()        // 1
				th.Throttle()        // 0
				th.SuccessfulRPC()   // +2 => 2
				return th.Throttle() // ok
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := tt.setup()
			got := tt.actions(th)
			if got != tt.want {
				t.Fatalf("got: %v want: %v", got, tt.want)
			}
		})
	}
}
