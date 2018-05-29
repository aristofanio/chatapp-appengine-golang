package member

import "uuid"

type Repository interface {
	Store(*Entity) error
	FindByID(uuid.UID) (*Entity, error)
	ListAllCreatedMsg() ([]*Entity, error)
	ListAllPendentMsg(uuid.UID) ([]*Entity, error)
}
