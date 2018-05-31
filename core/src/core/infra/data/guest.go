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

type Guest struct {
	ID                 uuid.UID `datastore:"id" json:"id"`
	DeviceUUID         string   `datastore:"device_uuid"`
	DeviceVersion      string   `datastore:"device_version"`
	DeviceSO           string   `datastore:"device_so"`
	DeviceManufacturer string   `datastore:"device_manufacturer"`
	DeviceSerial       string   `datastore:"device_serial"`
	CreatedIn          int64    `datastore:"created_in,noindex"`
	UpdatedIn          int64    `datastore:"updated_in,noindex"`
}

//------------------------------------------------------------------
// Guest Data Manager
//------------------------------------------------------------------

type GuestDataMgr struct {
	ctx context.Context
}

func (m GuestDataMgr) NewGuest() *Guest {
	//create instance
	gu := new(Guest)
	gu.ID = uuid.NewUID("guest")
	gu.CreatedIn = utils.Now()
	gu.UpdatedIn = utils.Now()
	//result
	return gu
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
