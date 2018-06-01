package guest

import "core/err"

var (
	ErrFailOnRegister = err.NewErr(3001, "Dispositivo não pode ser registrado")
)
