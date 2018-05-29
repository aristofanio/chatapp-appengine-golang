package message

import (
	"context"
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
	msgRep  Repository
	mbrRep  member.Repository
	ctx     context.Context
	notServ notifier.Service
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
	err := m.msgRep.Store(met)
	if err != nil {
		return "", err
	}
	//notify
	m.notServ.Fire(notifier.CreatedMsgEvt, notifier.Evt{met.Origin, met.Destiny, met.ID})
	//send notification for dst
	m.sendMsg(dst, mid)
	//result msg uuid
	return "uuid", nil
}

func (m Manager) findMember(userID uuid.UID) (*member.Entity, error) {
	//
	mbr, err := m.mbrRep.FindByID(userID)
	if err != nil {
		return mbr, err
	}
	return mbr, nil
}

func (m Manager) sendMsg(mbr member.Entity, msgID uuid.UID) {
	//notify member
	m.notServ.Fire(notifier.SentMsgEvt, notifier.Evt{m.user.ID, mbr.ID, msgID})
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
	msg, err := m.msgRep.FindByID(msgID)
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
	err = m.msgRep.Store(msg)
	if err != nil {
		return err
	}
	//fire event
	m.notServ.Fire(notifier.DestroyedMgOrgEvt, notifier.Evt{m.user.ID, msg.Destiny, msgID})
	//result
	return nil
}

func (m Manager) ListAllPendent(msgID uuid.UID) []Entity {
	//result
	result := make([]Entity, 0)
	//list all pendent
	msgs, err := m.msgRep.ListAllPendentMsg(m.user.ID)
	if err == nil {
		for _, msg := range msgs {
			result = append(result, *msg)
			//fire event
			m.notServ.Fire(notifier.ReadDstEvt, notifier.Evt{msg.Destiny, m.user.ID, msgID})
		}
	}
	//result
	return result
}

func (m Manager) ConfirmDestroingMsg(msgID uuid.UID) error {
	//retrieve message by ref (to msgID)
	msg, err := m.msgRep.FindByID(msgID)
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
	err = m.msgRep.Store(msg)
	if err != nil {
		return err
	}
	//fire event
	m.notServ.Fire(notifier.DestroyedMgDstEvt, notifier.Evt{msg.Destiny, m.user.ID, msgID})
	//result
	return nil
}

func (m Manager) SendAllMsg() { //each 1s by task auto - idempotent
	//list all created
	msgs, err := m.msgRep.ListAllCreatedMsg()
	if err != nil {
		return
	}
	//send for each e
	for _, msg := range msgs {
		mbr, err := m.findMember(msg.Destiny)
		if err != nil {
			continue
		}
		m.sendMsg(*mbr, msg.ID)
	}
}

func NewManager(msgRep Repository, mbrRep member.Repository, not notifier.Service,
	ctx context.Context, orig member.Entity) Manager {
	return Manager{msgRep, mbrRep, ctx, not, orig}
}
