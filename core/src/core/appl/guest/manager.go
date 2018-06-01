package guest

import (
	"context"
	"core/infra/data"
)

type Manager interface {
	Register(uuid, version, serial, os, mnftr string) (string, error)
}

type managerInst struct {
	ctx context.Context
}

func (m managerInst) Register(uuid, version, serial, os, mnftr string) (string, error) {
	//managers
	guestDataMgr := data.NewGuestDataMgr(m.ctx)
	guest := guestDataMgr.NewGuest(uuid, version, serial, os, mnftr)
	err := guestDataMgr.Store(guest)
	if err != nil {
		return "", ErrFailOnRegister.Original(err)
	}
	//result success
	return string(guest.ID), nil
}

func NewManager(ctx context.Context) Manager {
	return managerInst{ctx}
}
