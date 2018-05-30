package err

import (
	"fmt"
)

type GeneralErr struct {
	errCode     int
	errText     string
	errOriginal error
}

func (e *GeneralErr) Original(err error) *GeneralErr {
	e.errOriginal = err
	return e
}

func (e *GeneralErr) Error() string {
	msgErr := fmt.Sprintf("[ErrCode: %d] %s - Original: {%s}", e.errCode, e.errText, e.errOriginal.Error())
	if e.errOriginal == nil {
		msgErr = fmt.Sprintf("[ErrCode: %d] %s", e.errCode, e.errText)
	}
	return msgErr
}

func (e *GeneralErr) ErrorJson() string {
	msgErr := fmt.Sprintf("{\"code\": %d, \"text\": %s, \"original\": %s}", e.errCode, e.errText, e.errOriginal.Error())
	return msgErr
}

func NewErr(code int, text string) *GeneralErr {
	return &GeneralErr{code, text, nil}
}
