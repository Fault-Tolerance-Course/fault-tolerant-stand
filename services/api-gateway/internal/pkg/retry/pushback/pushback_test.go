package pushback

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestPushbackParse(t *testing.T) {
	tests := []struct {
		name string
		md   metadata.MD
		want Pushback
	}{
		{
			name: "nil metadata → no pushback",
			md:   nil,
			want: Pushback{},
		},
		{
			name: "no header → no pushback",
			md:   metadata.Pairs("other", "123"),
			want: Pushback{},
		},
		{
			name: "-1 → reject",
			md:   metadata.Pairs("grpc-retry-pushback-ms", "-1"),
			want: Pushback{Has: true, Reject: true},
		},
		{
			name: "delay = 150ms",
			md:   metadata.Pairs("grpc-retry-pushback-ms", "150"),
			want: Pushback{Has: true, Delay: 150 * time.Millisecond},
		},
		{
			name: "invalid number → ignore",
			md:   metadata.Pairs("grpc-retry-pushback-ms", "not-num"),
			want: Pushback{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.md)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got: %#v want: %#v", got, tt.want)
			}
		})
	}
}
