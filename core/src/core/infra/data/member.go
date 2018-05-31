package data

import (
	"context"
	"core/infra/data/uuid"
)

//------------------------------------------------------------------
// Member Data
//------------------------------------------------------------------

const memberKind = "Member"

type Member struct {
	ID        uuid.UID `datastore:"id" json:"id"`
	Name      string   `datastore:"name" json:"name"`
	Descr     string   `datastore:"descr" json:"descr"`
	FToken    string   `datastore:"ftoken" json:"ftoken"`
	PhotoURL  string   `datastore:"photo" json:"photo"`
	CreatedIn int64    `datastore:"created_in" json:"created_in"`
	UpdatedIn int64    `datastore:"updated_in" json:"updated_in"`
}

//------------------------------------------------------------------
// Member Data Manager
//------------------------------------------------------------------

type MemberDataMgr struct {
	ctx context.Context
}

func (mgr MemberDataMgr) Store(m *Member) error {
	return storeEntity(mgr.ctx, memberKind, m.ID, m)
}

func (mgr MemberDataMgr) Get(uid uuid.UID) (*Member, error) {
	//result
	rslt := &Member{}
	//get entity
	err := findEntity(mgr.ctx, memberKind, uid, rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func NewMemberDataMgr(ctx context.Context) MemberDataMgr {
	return MemberDataMgr{ctx}
}
