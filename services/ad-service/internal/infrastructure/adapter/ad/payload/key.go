package payload

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

const (
	Prefix = "advertisement"
)

func AdKey(adID uuid.UUID, version string) string {
	// {version}:{prefix}:{uuid}
	return fmt.Sprintf("%s:%s:%s", version, Prefix, adID.String())
}

func AdID(key string) uuid.UUID {
	parts := strings.Split(key, ":")
	if len(parts) != 3 {
		return uuid.Nil
	}
	return uuid.FromStringOrNil(parts[2])
}
