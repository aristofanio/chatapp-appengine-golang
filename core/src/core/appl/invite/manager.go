package invite

import (
	"context"
	"core/appl/notifier"
	"core/infra/data"
	"core/infra/data/uuid"
)

type Manager interface {

	//create an invite
	Create(toId string) (*data.Invite, error)

	//destroy an invite
	Destroy(toId string) error

	//accept an invite
	Accept(inviteID string) (*data.Chat, error)
}

type managerInst struct {
	usr  *data.User
	ctx  context.Context
	push notifier.Manager
}

//-------------------------------------------------------------------
// Methods - implementations of Entity
//-------------------------------------------------------------------

//TODO: criar apenas um convite para cada par e com o mesmo prazo de validade
func (m managerInst) Create(toId string) (*data.Invite, error) {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	invDataMgr := data.NewInviteDataMgr(m.ctx)
	//get other member
	otherID := uuid.UID(toId)
	other, err := mbrDataMgr.Get(otherID)
	if err != nil {
		return nil, err
	}
	//check invite to
	oldInviteTo := invDataMgr.GetByPair(m.usr.ID, other.ID)
	if oldInviteTo != nil {
		return nil, ErrAlreadyExistsinviteTo
	}
	oldInviteFrom := invDataMgr.GetByPair(other.ID, m.usr.ID)
	if oldInviteFrom != nil {
		return nil, ErrAlreadyExistsinviteFrom
	}
	//TODO: check blacklist
	//save invite
	invite := invDataMgr.NewInvite(m.usr.ID, other.ID)
	invDataMgr.Store(invite)
	//send invite
	m.push.SendInvite(other.ID)
	//result
	return invite, nil
}

func (m managerInst) Destroy(toId string) error {
	//data managers
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	invDataMgr := data.NewInviteDataMgr(m.ctx)
	//get other member
	otherID := uuid.UID(toId)
	other, err := mbrDataMgr.Get(otherID)
	if err != nil {
		return err
	}
	//check invite to
	oldInviteFrom := invDataMgr.GetByPair(other.ID, m.usr.ID)
	if oldInviteFrom != nil {
		return ErrAlreadyExistsinviteFrom
	}
	oldInviteFrom.Reject()
	//save invite and result
	return invDataMgr.Store(oldInviteFrom)
}

func (m managerInst) Accept(inviteID string) (*data.Chat, error) {
	//data managers
	invDataMgr := data.NewInviteDataMgr(m.ctx)
	mbrDataMgr := data.NewMemberDataMgr(m.ctx)
	//get invite
	invUID := uuid.UID(inviteID)
	invite, err := invDataMgr.Get(invUID)
	if err != nil {
		return nil, err
	}
	//get to and from members
	to, _ := mbrDataMgr.Get(invite.To)
	from, _ := mbrDataMgr.Get(invite.From)
	//accept chat and create invite
	chtDataMgr := data.NewChatDataMgr(m.ctx)
	chat, err := chtDataMgr.NewChat(invite.ID, *to, *from)
	if err != nil {
		return nil, err
	}
	invite.Accept()
	invDataMgr.Store(invite)
	//return chat (success)
	return chat, nil
}
