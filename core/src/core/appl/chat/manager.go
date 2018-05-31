package chat

type Manager interface {

	//create message
	CreateMessage()

	//destroy message
	DestroyMessage()

	//get messages for user
	GetMessages()

	//mark message as read
	MarkMessageAsRead()

	//mark message as received
	MarkMessageAsReceived()

	//mark message as destroyed
	MarkMessageAsDestroyed()
}
