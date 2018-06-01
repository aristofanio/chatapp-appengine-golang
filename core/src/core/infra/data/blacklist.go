package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"
	"fmt"

	"google.golang.org/appengine/datastore"
)

//-------------------------------------------------------------------
// Variables, constants and types (struct, interface, func)
//-------------------------------------------------------------------

const blackListKind = "BlackList"

type BlackList struct {
	ID        uuid.UID `datastore:"id"`
	Guest0    uuid.UID `datastore:"guest0"`
	Guest1    uuid.UID `datastore:"guest1"`
	IsAuto    bool     `datastore:"is_auto"`
	CreatedIn int64    `datastore:"created_in,noindex"`
	UpdatedIn int64    `datastore:"updated_in,noindex"`
}

//-------------------------------------------------------------------
// Operations
//-------------------------------------------------------------------

type BlackListDataMgr struct {
	ctx context.Context
}

func (m BlackListDataMgr) Store(e *BlackList) error {
	return storeEntity(m.ctx, blackListKind, e.ID, e)
}

func (m BlackListDataMgr) AddToBlackList(fromId, toId uuid.UID) error {
	//
	crr := utils.Now()
	//
	id0 := uuid.UID(fmt.Sprintf("%s:%s", fromId, toId))
	bl0 := &BlackList{ID: id0, Guest0: fromId, Guest1: toId,
		IsAuto: false, CreatedIn: crr, UpdatedIn: crr}
	err0 := m.Store(bl0)
	if err0 != nil {
		return err0
	}
	//
	id1 := uuid.UID(fmt.Sprintf("%s:%s", toId, fromId))
	bl1 := &BlackList{ID: id1, Guest0: toId, Guest1: fromId,
		IsAuto: true, CreatedIn: crr, UpdatedIn: crr}
	err1 := m.Store(bl1)
	if err1 != nil {
		return err1
	}
	//
	return nil
}

func (m BlackListDataMgr) IsBlackList(fromId, toId uuid.UID) bool {
	//
	id0 := fmt.Sprintf("%s:%s", fromId, toId)
	id1 := fmt.Sprintf("%s:%s", toId, fromId)
	//
	key0 := datastore.NewKey(m.ctx, blackListKind, id0, 0, nil)
	key1 := datastore.NewKey(m.ctx, blackListKind, id1, 0, nil)
	//
	var dst interface{}
	if err := datastore.GetMulti(m.ctx, []*datastore.Key{key0, key1}, dst); err != nil {
		return false
	}
	//
	return dst != nil

}

func NewBlackListDataMgr(ctx context.Context) BlackListDataMgr {
	return BlackListDataMgr{ctx}
}
