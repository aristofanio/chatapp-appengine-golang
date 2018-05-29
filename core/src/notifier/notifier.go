package notifier

import (
	"member"
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

type DataEvt interface {
}

type Service interface {
	//taskqueue.NewPOSTTask("")
	Fire(evtType int, org member.Entity, dt DataEvt) error
}
