package invite

import "core/err"

var (
	ErrNotPermissionOperation  = err.NewErr(5001, "Não há permissão para realizar esta operação")
	ErrAlreadyExistsinviteFrom = err.NewErr(5002, "Convite já existe. Verifique a aba de contatos")
	ErrAlreadyExistsinviteTo   = err.NewErr(5003, "Convite já existe. Aguarde a resposta")
)
