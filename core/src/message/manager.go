package message

import (
	"context"
	"data"
	"member"
	"notifier"
	"utils"
	"uuid"
)

const pfx = "msg"

type Manager struct {
	ctx      context.Context
	notifier notifier.Register
	manager  data.MessageManage
	origin   member.Entity
}

func (m Manager) CreateTextMsg(dst member.Entity, msg string) (string, error) {
	//create message
	mid := uuid.NewUID(pfx)
	now := utils.Now()
	met := Entity{
		ID:           mid,
		Origin:       m.origin.ID,
		Destiny:      dst.ID,
		Type:         TextMsgType,
		Text:         MsgText(msg),
		OriginState:  CreatedState,
		DestinyState: CreatedState,
		CreatedIn:    now,
		UpdatedIn:    now,
	}
	//store message
	//TODO: notify
	m.not.On(notifier.CreatedMsgEvt, dst)
	//result msg uuid
	return "uuid", nil
}

func (m Manager) CreateImgMsg(dst member.Entity, img string) (string, error) {
	m.not.On(notifier.CreatedMsgEvt, dst)
	return "uuid", nil
}

func (m Manager) CreateHybMsg(dst member.Entity, msg, img string) (string, error) {
	m.not.On(notifier.CreatedMsgEvt, dst)
	return "uuid", nil
}

func (m Manager) DestroyMsg(msgID string) error {
	//
	var dst member.Entity //TODO: recuperar + validar org = msg.Org
	//
	m.not.On(notifier.DestroyedMgOrgEvt, dst)
	//
	return nil
}

func (m Manager) ConfirmDestroingMsg(msgID string) error {
	//
	var dst member.Entity //TODO: recuperar + validar dst = msg.Dst
	//
	m.not.On(notifier.DestroyedMgOrgEvt, dst) //finaliza a notificaca
	//
	return nil
}

func (m Manager) SendMsg(msgID string) error {
	//todas as pendets -> criar notifications
	var dst member.Entity //TODO: recuperar + validar dst = msg.Dst
	//
	m.not.On(notifier.SentMsgEvt, dst) //finaliza a notificaca
	//
	return nil
}

func (m Manager) SendAllMsg() error { //each 1s by task auto
	//todas as pendets -> criar notifications
	return nil
}

func NewManager(ctx context.Context, not notifier.Register, orig member.Entity) Manager {
	return Manager{ctx, not, orig}
}
