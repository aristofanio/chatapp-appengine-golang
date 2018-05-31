package data

import "core/err"

//Errors
var (
	ErrEntityNotFound      = err.NewErr(1001, "Entidade não encontrada")
	ErrEntityNotDeleted    = err.NewErr(1002, "Entidade não pode ser excluída")
	ErrEntitiesNotListed   = err.NewErr(1003, "Entidades não podem ser listadas")
	ErrEntityMustBeNotNull = err.NewErr(1004, "Entidade não pode ser nula")
)
