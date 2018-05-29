package message

import (
	"context"
	"data"
	"errors"
	"member"
	"notifier"
	"utils"
	"uuid"
)

const pfx = "msg"

var (
	InvalidOriginErr = errors.New("Invalid origin error")
	InvalidMemberErr = errors.New("Invalid member error")
)

type Manager struct {
	ctx     context.Context
	notServ notifier.Service
	dataMgr data.Manager
	user    member.Entity
}

func (m Manager) createMsg(typ int, dst member.Entity, msg, img string) (string, error) {
	//create message
	mid := uuid.NewUID(pfx)
	now := utils.Now()
	met := &Entity{
		ID:           mid,
		Origin:       m.user.ID,
		Destiny:      dst.ID,
		Type:         typ,
		Text:         MsgText(msg),
		ImageURL:     MsgImage(msg),
		OriginState:  CreatedState,
		DestinyState: CreatedState,
		CreatedIn:    now,
		UpdatedIn:    now,
	}
	//store message
	err := m.dataMgr.Store(met)
	if err != nil {
		return "", err
	}
	//notify
	m.notServ.Fire(notifier.CreatedMsgEvt, dst, mid)
	//send notification for dst
	m.sendMsg(dst, mid)
	//result msg uuid
	return "uuid", nil
}

func (m Manager) findMember(uid uuid.UID) (member.Entity, error) {
	//
	mbr := member.Entity{}
	err := m.dataMgr.FindByID(&mbr, uid)
	if err != nil {
		return mbr, err
	}
	return mbr, nil
}

func (m Manager) sendMsg(mbr member.Entity, msgID uuid.UID) {
	//notify member
	m.notServ.Fire(notifier.SentMsgEvt, mbr, msgID)
}

func (m Manager) CreateTextMsg(dst member.Entity, msg string) (string, error) {
	//result msg uuid
	return m.createMsg(TextMsgType, dst, msg, "")
}

func (m Manager) CreateImgMsg(dst member.Entity, img string) (string, error) {
	//result msg uuid
	return m.createMsg(ImageMsgType, dst, "", img)
}

func (m Manager) CreateHybMsg(dst member.Entity, msg, img string) (string, error) {
	//result msg uuid
	return m.createMsg(HybridMsgType, dst, msg, img)
}

func (m Manager) DestroyMsg(msgID uuid.UID) error {
	//retrieve message by ref (to msgID)
	msg := Entity{}
	err := m.dataMgr.FindByID(&msg, msgID)
	if err != nil {
		return err
	}
	//retrieve origin
	org, err := m.findMember(msg.Origin)
	if err != nil {
		return err
	}
	//validate (org = msg.User)
	if !m.user.ID.Equals(org.ID) {
		return InvalidOriginErr
	}
	//update state for destroyed
	msg.OriginState = DestroyedState
	err = m.dataMgr.Store(&msg)
	if err != nil {
		return err
	}
	//fire event
	m.notServ.Fire(notifier.DestroyedMgOrgEvt, org, msgID)
	//result
	return nil
}

func (m Manager) ConfirmDestroingMsg(msgID uuid.UID) error {
	//retrieve message by ref (to msgID)
	msg := Entity{}
	err := m.dataMgr.FindByID(&msg, msgID)
	if err != nil {
		return err
	}
	//retrieve destiny
	dst, err := m.findMember(msg.Destiny)
	if err != nil {
		return err
	}
	//validate (dst = msg.User)
	if !m.user.ID.Equals(dst.ID) {
		return InvalidOriginErr
	}
	//update state for destroyed
	msg.OriginState = DestroyedState
	err = m.dataMgr.Store(&msg)
	if err != nil {
		return err
	}
	//fire event
	m.notServ.Fire(notifier.DestroyedMgDstEvt, dst, msgID)
	//result
	return nil
}

func (m Manager) SendAllMsg() { //each 1s by task auto - idempotent
	//list all created
	msgs := m.dataMgr.ListAllCreated()
	//send for each e
	for _, e := range msgs {
		msg := e.(*Entity)
		mbr, err := m.findMember(msg.Destiny)
		if err != nil {
			continue
		}
		m.sendMsg(mbr, msg.ID)
	}
}

func NewManager(ctx context.Context, not notifier.Service,
	dt data.Manager, orig member.Entity) Manager {
	return Manager{ctx, not, dt, orig}
}
