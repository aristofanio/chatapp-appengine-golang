package message

import (
	"context"
	"uuid"

	"google.golang.org/appengine/datastore"
)

const (
	TextMsgType = iota
	ImageMsgType
	HybridMsgType
)

const (
	CreatedState = iota
	PendentState
	DeliveredState
	ReadState
	DestroyedState
)

const kind = "Message"

type MsgText string

type MsgImage string

type Entity struct {
	ID           uuid.UID `datastore:"id" json:"id"`
	Origin       uuid.UID `datastore:"origin" json:"origin"`
	Destiny      uuid.UID `datastore:"destiny" json:"destiny"`
	Type         int      `datastore:"type" json:"type"`
	Text         MsgText  `datastore:"text" json:"text"`
	ImageURL     MsgImage `datastore:"imgurl" json:"imgurl"`
	OriginState  int      `datastore:"origin_state" json:"origin_state"`
	DestinyState int      `datastore:"destiny_state" json:"destiny_state"`
	CreatedIn    int64    `datastore:"created_in" json:"created_in"`
	UpdatedIn    int64    `datastore:"updated_in" json:"updated_in"`
}

func (e *Entity) Key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, kind, string(e.ID), 0, nil)
}

func (e *Entity) EName() string {
	return kind
}
