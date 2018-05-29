package member

import (
	"context"
	"uuid"

	"google.golang.org/appengine/datastore"
)

const kind = "Member"

type Entity struct {
	ID        uuid.UID `datastore:"id" json:"id"`
	Name      string   `datastore:"name" json:"name"`
	Descr     string   `datastore:"descr" json:"descr"`
	FToken    string   `datastore:"ftoken" json:"ftoken"`
	PhotoURL  string   `datastore:"photo" json:"photo"`
	CreatedIn int64    `datastore:"created_in" json:"created_in"`
	UpdatedIn int64    `datastore:"updated_in" json:"updated_in"`
}

func (e *Entity) Key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, kind, string(e.ID), 0, nil)
}

func (e *Entity) EName() string {
	return kind
}
