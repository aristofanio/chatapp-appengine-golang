package data

import (
	"core/infra/data/uuid"
	"strconv"

	"golang.org/x/net/context"
)

//------------------------------------------------------------------
// Message Data
//------------------------------------------------------------------

const (
	messageKind = "Message"
)

const (
	TypeTextMsg = iota
	TypeImageMsg
	TypeHybridMsg
)

const (
	StateCreated = iota
	StatePendent
	StateDelivered
	StateRead
	StateDestroyed
)

type MsgText string

type MsgImage string

type Message struct {
	ID           uuid.UID `datastore:"id" json:"id"`
	Origin       uuid.UID `datastore:"origin" json:"origin"`
	Destiny      uuid.UID `datastore:"destiny" json:"destiny"`
	Type         int      `datastore:"type" json:"type"`
	Text         MsgText  `datastore:"text" json:"text"`
	ImageURL     MsgImage `datastore:"imgurl" json:"imgurl"`
	OriginState  int      `datastore:"origin_state" json:"origin_state"`
	DestinyState int      `datastore:"destiny_state" json:"destiny_state"`
	CreatedIn    int64    `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn    int64    `datastore:"updated_in,noindex" json:"update_in"`
}

//------------------------------------------------------------------
// Message Data Manager
//------------------------------------------------------------------

type MessageDataMgr struct {
	ctx context.Context
}

func (mgr MessageDataMgr) Store(m *Message) error {
	return storeEntity(mgr.ctx, memberKind, m.ID, m)
}

func (mgr MessageDataMgr) Get(uid uuid.UID) (*Message, error) {
	//result
	rslt := &Message{}
	//get entity
	err := findEntityByID(mgr.ctx, messageKind, uid, rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func (mgr MessageDataMgr) ListAllPendentsMsgs() ([]*Message, error) {
	//result
	rslts := make([]*Message, 0)
	//filter
	filter := make(map[string]string)
	filter["origin_state="] = strconv.Itoa(StatePendent)
	//list entities
	err := listEntities(mgr.ctx, messageKind, filter, "", rslts)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslts, nil
}

func NewMessageDataMgr(ctx context.Context) MessageDataMgr {
	return MessageDataMgr{ctx}
}
