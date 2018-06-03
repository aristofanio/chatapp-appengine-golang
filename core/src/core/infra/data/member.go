package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"
)

//------------------------------------------------------------------
// Member Data
//------------------------------------------------------------------

const memberKind = "Member"

type Member struct {
	ID        uuid.UID `datastore:"mid" json:"mid"`
	Nick      string   `datastore:"nick" json:"nick"`
	Name      string   `datastore:"name" json:"name"`
	Email     string   `datastore:"email" json:"email"`
	Gender    string   `datastore:"gender" json:"gender"`
	Photo     string   `datastore:"photo" json:"photo"`
	Age       int      `datastore:"age" json:"age"`
	IsRemoved bool     `datastore:"is_removed" json:"is_removed"`
	IsBlocked bool     `datastore:"is_blocked" json:"is_blocked"`
	CreatedIn int64    `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn int64    `datastore:"updated_in,noindex" json:"update_in"`
}

//------------------------------------------------------------------
// Member Data Manager
//------------------------------------------------------------------

type MemberDataMgr struct {
	ctx context.Context
}

func (m MemberDataMgr) NewMember(name, nick, email string) *Member {
	//create instance
	u := new(Member)
	u.ID = uuid.NewUID("mbr")
	u.Nick = nick
	u.Name = name
	u.Email = email
	u.CreatedIn = utils.Now()
	u.UpdatedIn = utils.Now()
	//result
	return u
}

func (mgr MemberDataMgr) Store(m *Member) error {
	return storeEntity(mgr.ctx, memberKind, m.ID, m)
}

func (mgr MemberDataMgr) StoreWithUser(m *Member, u *User) error {
	return storeEntityWithParent(mgr.ctx,
		memberKind, m.ID, m,
		userKind, u.ID, u)
}

func (mgr MemberDataMgr) Remove(m *Member) error {
	return deleteEntity(mgr.ctx, memberKind, m.ID)
}

func (mgr MemberDataMgr) Get(uid uuid.UID) (*Member, error) {
	//result
	rslt := &Member{}
	//get entity
	err := findEntityByID(mgr.ctx, memberKind, uid, rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func (mgr MemberDataMgr) ListAll(offset, limit int) ([]*Member, error) {
	//result
	rslt := make([]*Member, 0)
	//filters
	filters := make(map[string]string)
	filters["is_blocked="] = "false"
	filters["is_removed="] = "false"
	//get entity
	err := listEntities(mgr.ctx, memberKind, filters, "", rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func NewMemberDataMgr(ctx context.Context) MemberDataMgr {
	return MemberDataMgr{ctx}
}
