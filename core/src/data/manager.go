package data

import (
	"context"

	"google.golang.org/appengine/datastore"
)

type Entity interface {
	Key(context.Context) *datastore.Key
	EName() string
}

type MessageManager interface {
	Store(Entity) error
	FindByID(uid string) Entity
}
