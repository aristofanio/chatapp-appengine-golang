package comm

import (
	"core/err"
)

var (
	ErrFailOnSend = err.NewErr(1011, "Fail on send fcm notification")
)
