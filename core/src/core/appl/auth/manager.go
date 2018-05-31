package auth

import (
	"context"
	"core/infra/data"
	"core/infra/data/uuid"
	"errors"
)

//Manager manages the authetication
type Manager interface {

	//Auth authenticates user creating a session and returning session token
	Auth(uGuestId, uEmail, uPass, uFcmToken string, lat, lng float64) (*data.Access, error)

	//DeAuth deauthenticates user by session token
	DeAuth(sToken string) error

	//CheckAuth checks the validity of session token
	CheckAuth(sToken string) (bool, error)
}

//An instance of Manager
type managerInst struct {
	ctx context.Context
}

func (m managerInst) checkGuest(uGuestId string) (*data.Guest, error) {
	//managers
	guestDataMgr := data.NewGuestDataMgr(m.ctx)
	//find guest (created on first open app)
	guestUID := uuid.UID(uGuestId)
	guest := guestDataMgr.GetGuest(guestUID)
	if guest == nil {
		return nil, errors.New("Requisição para usuário bloqueado ou inexistente.")
	}
	return guest, nil
}

func (m managerInst) Auth(uGuestId, uEmail, uPass, uFcmToken string, lat, lng float64) (*data.Access, error) {
	//managers
	accessDataMgr := data.NewAccessDataMgr(m.ctx)
	userDataMgr := data.NewUserDataMgr(m.ctx)
	//find guest (created on first open app)
	guest, err := m.checkGuest(uGuestId)
	if err != nil {
		return nil, err
	}
	//retrieve user by nick
	user, err := userDataMgr.GetUserByEmail(uEmail)
	if err != nil {
		return nil, err
	}
	//check password
	if passed := user.CheckPassword(uPass); passed {
		return nil, ErrAuthNotAccepted
	}
	//invalidar todos os acessos para o dispositivo
	err = accessDataMgr.InvalidAllOldAccess(*guest)
	if err != nil {
		return nil, err
	}
	//check position
	if lat == 0 || lng == 0 {
		return nil, ErrGPSDataInvalid.Original(err)
	}
	//create session
	return accessDataMgr.NewAccess(*guest, *user, uFcmToken, lat, lng)
}

func (m managerInst) DeAuth(uGuestId, sToken string) error {
	//managers
	accessDataMgr := data.NewAccessDataMgr(m.ctx)
	//find guest (created on first open app)
	guest, err := m.checkGuest(uGuestId)
	if err != nil {
		return err
	}
	//retrieve access
	access, err := accessDataMgr.GetAccess(uuid.UID(sToken))
	if err != nil {
		return err
	}
	//return success
	if access.GuestID.Equals(guest.ID) {
		return accessDataMgr.InvalidAccess(access)
	}
	//return fail
	return ErrNotPermissionOperation
}

func (m managerInst) CheckAuth(sToken string) (bool, error) {
	//managers
	accessDataMgr := data.NewAccessDataMgr(m.ctx)
	//retrieve access
	access, err := accessDataMgr.GetAccess(uuid.UID(sToken))
	if err != nil {
		return false, err
	}
	//return valid
	if access.IsValid {
		return true, nil
	}
	//return invalid
	return false, nil
}
