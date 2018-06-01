package comm

import (
	"net/http"

	"github.com/aristofanio/go-fcm"
)

type Pair struct {
	Key   string
	Value string
}

type FCMResp struct {
	Code    int
	Success bool
	ErrMsg  string
}

type FCMNotifier interface {
	SendInvite()
	AcceptInvite()
	BlockMember()
	SendMessage()
	ConfirmReadMessage()
	ConfirmReceivedMessage()
	ConfirmDestroyedMessage()
	//
	PushMessage(ftoken string, message string)
	PushData(ftoken string, data map[string]string)
}

type fcmNotifierInst struct {
	serverKey string
	client    *http.Client
}

func (f fcmNotifierInst) push(ftoken string, data interface{}, payload *fcm.NotificationPayload) FCMResp {
	//
	fcmClient := fcm.NewClient(f.serverKey)
	fcmClient.SetHTTPClient(f.client)
	//set priority
	fcmClient.SetPriority(fcm.HighPriority)
	//set notification message
	fcmClient.SetNotification(payload)
	//send message
	fcmClient.PushSingle(ftoken, data)
	//fcmClient.SetContentAvailable(true)
	fresp, err := fcmClient.Send()
	if err != nil {
		errCustom := ErrFailOnSend.Original(err)
		return FCMResp{
			Code:    errCustom.GetCode(),
			ErrMsg:  errCustom.GetText(),
			Success: false,
		}
	}
	//
	return FCMResp{
		Code:    fresp.StatusCode,
		ErrMsg:  fresp.Err,
		Success: fresp.Success == 1,
	}
}

func (f fcmNotifierInst) PushData(ftoken, title, msg, icon string, pairs []Pair) FCMResp {
	// create data
	data := make(map[string]string)
	for _, p := range pairs {
		data[p.Key] = p.Value
	}
	// message for background
	note := &fcm.NotificationPayload{
		Title:       title,
		Body:        msg,
		Icon:        icon,
		ClickAction: "FCM_PLUGIN_ACTIVITY",
		Sound:       "default"}
	//push
	return f.push(ftoken, data, note)
}

func (f fcmNotifierInst) PushMessage(ftoken, title, msg, icon string) FCMResp {
	//create notification message
	note := &fcm.NotificationPayload{
		Title: title,
		Body:  msg,
		Icon:  icon,
		Sound: "default"}
	//push
	return f.push(ftoken, nil, note)
}
