package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"
)

//------------------------------------------------------------------
// Guest Data
//------------------------------------------------------------------

const guestKind = "Guest"

type Device struct {
	UUID         string `datastore:"uuid"`
	Version      string `datastore:"version"`
	OS           string `datastore:"os"`
	Manufacturer string `datastore:"manufacturer"`
	Serial       string `datastore:"serial"`
}

type Guest struct {
	ID        uuid.UID `datastore:"id" json:"id"`
	Device    Device   `datastore:"device" json:"device"`
	CreatedIn int64    `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn int64    `datastore:"updated_in,noindex" json:"update_in"`
}

//------------------------------------------------------------------
// Guest Data Manager
//------------------------------------------------------------------

type GuestDataMgr struct {
	ctx context.Context
}

func (m GuestDataMgr) NewGuest(did, version, serial, os, mnftr string) *Guest {
	//
	device := Device{
		UUID:         did,
		Version:      version,
		OS:           os,
		Manufacturer: mnftr,
		Serial:       serial,
	}
	//
	guest := &Guest{
		ID:        uuid.NewUID("guest"),
		Device:    device,
		CreatedIn: utils.Now(),
		UpdatedIn: utils.Now(),
	}
	//result
	return guest
}

func (m GuestDataMgr) GetGuest(uid uuid.UID) *Guest {
	gu := new(Guest)
	if err := findEntityByID(m.ctx, guestKind, uid, gu); err != nil {
		return nil
	}
	return gu
}

func (m GuestDataMgr) Store(gu *Guest) error {
	return storeEntity(m.ctx, guestKind, gu.ID, gu)
}

func NewGuestDataMgr(ctx context.Context) GuestDataMgr {
	return GuestDataMgr{ctx}
}
