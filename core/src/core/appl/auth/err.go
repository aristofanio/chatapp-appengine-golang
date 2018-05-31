package auth

import "core/err"

var (
	ErrAuthNotAccepted        = err.NewErr(2001, "Usuário e senha não combinam")
	ErrGPSDataInvalid         = err.NewErr(2002, "Dados do GPS inválidos. Por favor verifique as permissões ou reinicie a aplicação e tente novamente")
	ErrAccessNotFound         = err.NewErr(2003, "Acesso não encontrado")
	ErrNotPermissionOperation = err.NewErr(2004, "Não há permissão para realizar esta operação")
)
