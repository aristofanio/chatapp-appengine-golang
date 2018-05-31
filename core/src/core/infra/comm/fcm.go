package comm

import (
	"core/err"
	"net/http"

	"github.com/aristofanio/go-fcm"
)

var (
	ErrFailOnSend = err.NewErr(1003, "Fail on send fcm notification")
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

type FCMNotifierInst struct {
	serverKey string
	client    *http.Client
}

func (f FCMNotifierInst) push(ftoken string, data interface{}, payload *fcm.NotificationPayload) FCMResp {
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

func (f FCMNotifierInst) PushData(ftoken string, pairs []Pair) FCMResp {
	// create data
	data := make(map[string]string)
	for _, p := range pairs {
		data[p.Key] = p.Value
	}
	//push
	return f.push(ftoken, data, nil)
}

func (f FCMNotifierInst) PushMessage(ftoken, title, msg string) FCMResp {
	//create notification message
	note := &fcm.NotificationPayload{
		Title: title,
		Body:  msg,
		Sound: "default"}
	//push
	return f.push(ftoken, nil, note)
}
