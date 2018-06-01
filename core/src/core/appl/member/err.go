package member

import "core/err"

var (
	ErrNotPermissionOperation = err.NewErr(4001, "Não há permissão para realizar esta operação")
	ErrNotAllowSelfBlocking   = err.NewErr(4002, "Não há permissão para auto bloqueio")
)
