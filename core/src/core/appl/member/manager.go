package member

import (
	"context"
	"core/infra/data"
	"core/infra/data/uuid"
)

type Manager interface {

	//check nick name
	//TODO: resolver
	//CheckNick(nick string) bool

	//create a member
	Create(name, nick, email, password string) (*data.Member, error)

	//destroy a member
	Destroy(mid string) error

	//update a member
	UpdatePersonalInfo(mid, name, gender string, age int) error
	UpdatePhoto(mid, photo string) error
	UpdatePassword(mid, oldpassword, password string) error
	ResetPassword(mid, passtoken, password string) error

	//valid a member
	Verify(mid, passtoken string) error

	//block an other member
	BlockOther(otherMID string) error

	//list others
	ListOthers(offset, limit int) ([]data.Member, error)
}

type managerInst struct {
	usr *data.User
	ctx context.Context
}

func (m managerInst) Create(name, nick, email, password string) (*data.Member, error) {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	usrDataMgr := data.NewUserDataMgr(m.ctx)
	//create member
	member := mbrDataMgr.NewMember(name, email)
	//create user
	user := usrDataMgr.NewUser(member.ID, nick, email, password)
	//store member and user
	//TODO: in transaction
	//TODO: in valid email and nick like unique
	err := mbrDataMgr.Store(member)
	if err != nil {
		return nil, err
	}
	err = usrDataMgr.Store(user)
	if err != nil {
		mbrDataMgr.Remove(member)
		return nil, err
	}
	return member, nil
}

func (m managerInst) Destroy() error {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	//get member
	member, err := mbrDataMgr.Get(m.usr.ID)
	if err != nil {
		return err
	}
	member.IsRemoved = true
	//store member
	err = mbrDataMgr.Store(member)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) UpdatePersonalInfo(name, gender string, age int) error {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	//get member
	member, err := mbrDataMgr.Get(m.usr.ID)
	if err != nil {
		return err
	}
	member.Age = age
	member.Gender = gender
	member.Name = name
	//store member
	err = mbrDataMgr.Store(member)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) UpdatePhoto(photo string) error {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	//get member
	member, err := mbrDataMgr.Get(m.usr.ID)
	if err != nil {
		return err
	}
	//update member
	member.Photo = photo
	//update user
	m.usr.Photo = photo
	//store member (in transaction)
	err = mbrDataMgr.StoreWithUser(member, m.usr)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) UpdatePassword(oldpassword, password string) error {
	//data managers
	usrDataMgr := data.NewUserDataMgr(m.ctx)
	//update password
	match := m.usr.CheckPassword(oldpassword)
	if !match {
		return ErrNotPermissionOperation
	}
	m.usr.SetPassKey(password)
	//store user
	err := usrDataMgr.Store(m.usr)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) ResetPassword(passtoken, password string) error {
	//data managers
	usrDataMgr := data.NewUserDataMgr(m.ctx)
	//update password
	match := m.usr.CheckPassToken(passtoken)
	if !match {
		return ErrNotPermissionOperation
	}
	m.usr.SetPassKey(password)
	//store user
	err := usrDataMgr.Store(m.usr)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) Verify(passtoken string) error {
	//data managers
	usrDataMgr := data.NewUserDataMgr(m.ctx)
	//update password
	match := m.usr.CheckPassToken(passtoken)
	if !match {
		return ErrNotPermissionOperation
	}
	m.usr.IsVerified = true
	//store user
	err := usrDataMgr.Store(m.usr)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) BlockOther(otherMID string) error {
	//other uid
	otherUID := uuid.UID(otherMID)
	//check
	if m.usr.ID.Equals(otherUID) {
		return ErrNotAllowSelfBlocking
	}
	//data managers
	usrDataMgr := data.NewUserDataMgr(m.ctx)
	blkDataMgr := data.NewBlackListDataMgr(m.ctx)
	//get other user
	other, err := usrDataMgr.Get(otherUID)
	if err != nil {
		return err
	}
	//todo: add to backlist
	err = blkDataMgr.AddToBlackList(m.usr.ID, other.ID)
	if err != nil {
		return err
	}
	//success
	return nil
}

func (m managerInst) ListOthers(offset, limit int) ([]data.Member, error) {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	//listing
	mbrs, err := mbrDataMgr.ListAll(offset, limit)
	if err != nil {
		return nil, err
	}
	//result
	rslts := make([]data.Member, 0)
	for _, mbr := range mbrs {
		//TODO: remove blacklist and me
		if mbr.ID.Equals(m.usr.ID) {
			rslts = append(rslts, *mbr)
		}
	}
	return rslts, nil
}
