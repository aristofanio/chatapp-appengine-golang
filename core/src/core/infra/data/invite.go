package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"

	"google.golang.org/appengine/datastore"
)

//-------------------------------------------------------------------
// Invite Data
//-------------------------------------------------------------------

const inviteKind = "Invite"

type Invite struct {
	ID         uuid.UID `datastore:"id" json:"id"`
	From       uuid.UID `datastore:"from" json:"from"`
	To         uuid.UID `datastore:"to" json:"to"`
	IsAccepted bool     `datastore:"is_accepted" json:"is_accepted"`
	IsRejected bool     `datastore:"is_rejected" json:"is_rejected"`
	CreatedIn  int64    `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn  int64    `datastore:"updated_in,noindex" json:"update_in"`
}

func (u *Invite) Reject() {
	u.IsRejected = true
	u.UpdatedIn = utils.Now()
}

func (u *Invite) Accept() {
	u.IsAccepted = true
	u.UpdatedIn = utils.Now()
}

//-------------------------------------------------------------------
// Operations - private
//-------------------------------------------------------------------

// func newInviteId() string {
// 	u := uuid.NewV1()
// 	return u.String()
// }

//-------------------------------------------------------------------
// Invite Manager Data
//-------------------------------------------------------------------

type InviteDataMgr struct {
	ctx context.Context
}

func (m InviteDataMgr) NewInvite(from, to uuid.UID) *Invite {
	//create instance
	invite := new(Invite)
	invite.ID = uuid.NewUID("invite")
	invite.From = from
	invite.To = to
	invite.IsRejected = false
	invite.IsAccepted = false
	invite.CreatedIn = utils.Now()
	invite.UpdatedIn = utils.Now()
	//result
	return invite
}

func (mgr InviteDataMgr) Store(e *Invite) error {
	return storeEntity(mgr.ctx, inviteKind, e.ID, e)
}

func (m InviteDataMgr) Get(uid uuid.UID) (*Invite, error) {
	e := new(Invite)
	if err := findEntityByID(m.ctx, inviteKind, uid, e); err != nil {
		return nil, ErrEntityNotFound.Original(err)
	}
	return e, nil
}

func (m InviteDataMgr) GetByPair(fromId, toId uuid.UID) *Invite {
	//result
	dst := make([]*Invite, 0)
	//query
	var q = datastore.NewQuery(inviteKind).
		Filter("from=", toId).
		Filter("to=", fromId).
		Filter("is_rejected=", false).
		Filter("is_accepted=", false).
		Limit(1)
	//get all
	if _, err := q.GetAll(m.ctx, &dst); err != nil {
		return nil
	}
	//return
	if len(dst) == 0 {
		return nil
	} else {
		return dst[0]
	}
}

func (m InviteDataMgr) ListAllFrom(uid uuid.UID) ([]Invite, error) {
	//result
	dst := make([]Invite, 0)
	//query
	var q = datastore.NewQuery(inviteKind).
		Filter("is_rejected=", false).
		Filter("is_accepted=", false).
		Filter("from=", string(uid)).
		Limit(rowCountsMax)
	//get all
	if _, err := q.GetAll(m.ctx, dst); err != nil {
		return nil, ErrEntitiesNotListed.Original(err)
	}
	//return
	return dst, nil
}

func (m InviteDataMgr) ListAllTo(uid uuid.UID) ([]Invite, error) {
	//result
	dst := make([]Invite, 0)
	//query
	var q = datastore.NewQuery(inviteKind).
		Filter("is_rejected=", false).
		Filter("is_accepted=", false).
		Filter("to=", string(uid)).
		Limit(rowCountsMax)
	//get all
	if _, err := q.GetAll(m.ctx, dst); err != nil {
		return nil, ErrEntitiesNotListed.Original(err)
	}
	//return
	return dst, nil
}

func NewInviteDataMgr(ctx context.Context) InviteDataMgr {
	return InviteDataMgr{ctx}
}
