package notifier

import (
	"uuid"
)

const (
	CreatedMsgEvt     = iota //origin - on create message
	SentMsgEvt               //auto - on sent to firebase
	DestroyedMgDstEvt        //destiny - on destroy message from destiny confirmation
	DestroyedMgOrgEvt        //origin - on destroy message from origin
	NotifiedDstEvt           //destiny - on to get messages (from destiny)
	NotifiedOrgEvt           //origin - on confirm NotifiedDstEvt (from origin)
	ReadDstEvt               //destiny - on confirm like read (from destiny)
	ReadOrgEvt               //origin - on confirm like read (from origin)
)

type Evt struct {
	Target uuid.UID
	Origin uuid.UID
	Data   uuid.UID
}

type Listener func(Evt)

type Service interface {
	Fire(int, Evt) error
	On(int, Listener) error
}

//taskqueue.NewPOSTTask("")
