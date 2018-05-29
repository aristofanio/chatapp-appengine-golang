package utils

import (
	"time"
)

func Now() int64 {
	return time.Now().UnixNano()
}
