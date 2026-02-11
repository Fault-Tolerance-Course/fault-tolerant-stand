package leader

import (
	"fmt"
	"os"
)

func instanceID() string {
	host, _ := os.Hostname()
	return fmt.Sprintf("%s-%d", host, os.Getpid())
}
