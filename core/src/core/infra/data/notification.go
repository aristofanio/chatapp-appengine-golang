package data

import (
	"context"
	"core/infra/data/uuid"
	"core/utils"
)

//------------------------------------------------------------------
// Notification Data
//------------------------------------------------------------------

const (
	notificationKind = "Notification"
)

type FCMResp struct {
	Code    int    `datastore:"code" json:"code"`
	Success bool   `datastore:"success" json:"success"`
	FailMsg string `datastore:"fail_msg" json:"fail_msg"`
}

type FCMReq struct {
	FToken  string `datastore:"ftoken" json:"ftoken"`
	Payload string `datastore:"payload" json:"payload"`
}

type Notification struct {
	ID        uuid.UID `datastore:"id" json:"id"`
	Req       FCMReq   `datastore:"fcm_req" json:"fcm_req"`
	Resp      FCMResp  `datastore:"fcm_resp" json:"fcm_resp"`
	CreatedIn int64    `datastore:"created_in,noindex" json:"create_in"`
	UpdatedIn int64    `datastore:"updated_in,noindex" json:"update_in"`
}

//------------------------------------------------------------------
// Notification Data Manager
//------------------------------------------------------------------

type NotificationDataMgr struct {
	ctx context.Context
}

func (mgr NotificationDataMgr) NewNotification(ftoken, jsonData string) *Notification {
	not := &Notification{}
	not.ID = uuid.NewUID("not")
	not.Req = FCMReq{ftoken, jsonData}
	not.CreatedIn = utils.Now()
	not.UpdatedIn = utils.Now()
	return not
}

func (mgr NotificationDataMgr) Store(not *Notification) error {
	return storeEntity(mgr.ctx, memberKind, not.ID, not)
}

func (mgr NotificationDataMgr) Get(uid uuid.UID) (*Notification, error) {
	//result
	rslt := &Notification{}
	//get entity
	err := findEntityByID(mgr.ctx, notificationKind, uid, rslt)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return rslt, nil
}

func NewNotificationDataMgr(ctx context.Context) NotificationDataMgr {
	return NotificationDataMgr{ctx}
}
