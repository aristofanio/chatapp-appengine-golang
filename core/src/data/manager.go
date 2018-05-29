package data

import (
	"context"
	"uuid"

	"google.golang.org/appengine/datastore"
)

type Entity interface {
	Key(context.Context) *datastore.Key
	EName() string
}

type Manager interface {
	Store(Entity) error
	FindByID(Entity, uuid.UID) error
	ListAllCreatedMsg() []Entity
	ListAllPendentMsg(uuid.UID) []Entity
}
