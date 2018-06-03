package notifier

import "core/infra/data/uuid"

type Manager interface {

	//send invite
	SendInvite(toUID uuid.UID)

	// AcceptInvite()
	// BlockMember()
	// SendMessage()
	// ConfirmReadMessage()
	// ConfirmReceivedMessage()
	// ConfirmDestroyedMessage()
}

// type managerInst struct {
// 	usr          *data.User
// 	ctx          context.Context
// 	fcmServerKey string
// 	httpClient   *http.Client
// }

// func (m managerInst) SendInvite(toUID uuid.UID) {
// 	//get token TODO: recuperar o token
// 	ftoken := ""
// 	//send invite
// 	notifier := comm.NewFCMNotifier(m.fcmServerKey, m.httpClient)
// 	notifier.PushMessage(ftoken, "Tô Solteiro", "Há convites aguardando sua avaliação.", "icon")
// }
