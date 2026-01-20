package negative

import (
	"time"

	"ad-service/config"
	"ad-service/internal/pkg/aderror"
)

type Payload struct {
	Code           int    `json:"code"`
	Classification string `json:"classification"`
	Message        string `json:"message"`
	BaseMessage    string `json:"base_message"`
}

func Classify(err error) (*aderror.Error, time.Duration, bool) {
	casted, ok := aderror.As(err)

	switch {
	case ok && aderror.IsCode(err, aderror.NotFound):
		return casted, config.Instance().Cache.Negative.NotFoundTTL, true
	default:
		return nil, 0, false
	}
}
