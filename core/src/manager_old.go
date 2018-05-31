package appl

// import (
// 	"context"
// 	"core/data"
// 	"core/data/uuid"
// 	"core/event"
// 	"core/utils"
// 	"errors"
// )

// const pfx = "msg"

// var (
// 	InvalidOriginErr = errors.New("Invalid origin error")
// 	InvalidMemberErr = errors.New("Invalid member error")
// )

// type Manager struct {
// 	user    data.Member
// 	notServ event.Bus
// 	msgRep  data.MessageDataMgr
// 	mbrRep  data.MemberDataMgr
// 	notRep  data.NotificationDataMgr
// }

// func (m Manager) createMsg(typ int, dst data.Member, msg, img string) (string, error) {
// 	//create message
// 	mid := uuid.NewUID(pfx)
// 	now := utils.Now()
// 	met := &data.Message{
// 		ID:           mid,
// 		Origin:       m.user.ID,
// 		Destiny:      dst.ID,
// 		Type:         typ,
// 		Text:         data.MsgText(msg),
// 		ImageURL:     data.MsgImage(msg),
// 		OriginState:  data.StateCreated,
// 		DestinyState: data.StateCreated,
// 		CreatedIn:    now,
// 		UpdatedIn:    now,
// 	}
// 	//store message
// 	err := m.msgRep.Store(met)
// 	if err != nil {
// 		return "", err
// 	}
// 	//notify
// 	m.notServ.Fire(event.CreatedMsgEvt, event.Evt{met.Origin, met.Destiny, met.ID})
// 	//send notification for dst
// 	m.sendMsg(dst, mid)
// 	//result msg uuid
// 	return "uuid", nil
// }

// func (m Manager) findMember(userID uuid.UID) (*data.Member, error) {
// 	//
// 	mbr, err := m.mbrRep.Get(userID)
// 	if err != nil {
// 		return mbr, err
// 	}
// 	return mbr, nil
// }

// func (m Manager) sendMsg(mbr data.Member, msgID uuid.UID) {
// 	//notify member
// 	m.notServ.Fire(event.SentMsgEvt, event.Evt{m.user.ID, mbr.ID, msgID})
// }

// func (m Manager) CreateTextMsg(dst data.Member, msg string) (string, error) {
// 	//result msg uuid
// 	return m.createMsg(data.TypeTextMsg, dst, msg, "")
// }

// func (m Manager) CreateImgMsg(dst data.Member, img string) (string, error) {
// 	//result msg uuid
// 	return m.createMsg(data.TypeImageMsg, dst, "", img)
// }

// func (m Manager) CreateHybMsg(dst data.Member, msg, img string) (string, error) {
// 	//result msg uuid
// 	return m.createMsg(data.TypeHybridMsg, dst, msg, img)
// }

// func (m Manager) DestroyMsg(msgID uuid.UID) error {
// 	//retrieve message by ref (to msgID)
// 	msg, err := m.msgRep.Get(msgID)
// 	if err != nil {
// 		return err
// 	}
// 	//retrieve origin
// 	org, err := m.findMember(msg.Origin)
// 	if err != nil {
// 		return err
// 	}
// 	//validate (org = msg.User)
// 	if !m.user.ID.Equals(org.ID) {
// 		return InvalidOriginErr
// 	}
// 	//update state for destroyed
// 	msg.OriginState = data.StateDestroyed
// 	err = m.msgRep.Store(msg)
// 	if err != nil {
// 		return err
// 	}
// 	//fire event
// 	m.notServ.Fire(event.DestroyedMgOrgEvt, event.Evt{m.user.ID, msg.Destiny, msgID})
// 	//result
// 	return nil
// }

// func (m Manager) ConfirmDestroingMsg(msgID uuid.UID) error {
// 	//retrieve message by ref (to msgID)
// 	msg, err := m.msgRep.Get(msgID)
// 	if err != nil {
// 		return err
// 	}
// 	//retrieve destiny
// 	dst, err := m.findMember(msg.Destiny)
// 	if err != nil {
// 		return err
// 	}
// 	//validate (dst = msg.User)
// 	if !m.user.ID.Equals(dst.ID) {
// 		return InvalidOriginErr
// 	}
// 	//update state for destroyed
// 	msg.DestinyState = data.StateDestroyed
// 	err = m.msgRep.Store(msg)
// 	if err != nil {
// 		return err
// 	}
// 	//fire event
// 	m.notServ.Fire(event.DestroyedMgDstEvt, event.Evt{msg.Destiny, m.user.ID, msgID})
// 	//result
// 	return nil
// }

// func (m Manager) ListAllPendentsMsgs() []*data.Message {
// 	//list all pendent
// 	msgs, err := m.msgRep.ListAllPendentsMsgs()
// 	if err != nil {
// 		return nil //FIXME: isto mascara um erro de fato?
// 	}
// 	//result
// 	return msgs
// }

// func (m Manager) CreateNotifications() {
// 	//list all created
// 	if msgs := m.ListAllPendentsMsgs(); msgs != nil {
// 		not := m.notRep.NewNotification(ftoken)
// 		//fire event
// 		m.notServ.Fire(event.ReadDstEvt, event.Evt{msg.Destiny, m.user.ID, msgID})
// 	}
// }

// func (m Manager) SendAllMsg() { //each 1s by task auto - idempotent
// 	//list all created
// 	msgs, err := m.msgRep.ListAllCreatedMsg()
// 	if err != nil {
// 		return
// 	}
// 	//send for each e
// 	for _, msg := range msgs {
// 		mbr, err := m.findMember(msg.Destiny)
// 		if err != nil {
// 			continue
// 		}
// 		m.sendMsg(*mbr, msg.ID)
// 	}
// }

// func NewManager(ctx context.Context, orig data.Member) Manager {
// 	mgr := Manager{}
// 	mgr.mbrRep = data.NewMemberDataMgr(ctx)
// 	mgr.msgRep = data.NewMessageDataMgr(ctx)
// 	mgr.notRep = data.NewNotificationDataMgr(ctx)
// 	return mgr
// }
