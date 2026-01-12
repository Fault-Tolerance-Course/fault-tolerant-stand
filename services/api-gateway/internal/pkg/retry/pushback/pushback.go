package pushback

import (
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
)

type Pushback struct {
	Delay  time.Duration
	Reject bool
	Has    bool
}

func Parse(md metadata.MD) Pushback {
	if md == nil {
		return Pushback{}
	}

	pairs := md.Get("grpc-retry-pushback-ms")
	if len(pairs) == 0 {
		return Pushback{}
	}

	n, err := strconv.Atoi(pairs[0])
	if err != nil {
		// если сервер прислал мусор — игнорируем
		return Pushback{}
	}

	pb := Pushback{Has: true}

	if n < 0 {
		pb.Reject = true
		return pb
	}

	pb.Delay = time.Duration(n) * time.Millisecond
	return pb
}
