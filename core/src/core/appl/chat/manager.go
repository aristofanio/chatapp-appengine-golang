package chat

type Manager interface {

	//create message
	CreateMessage()

	//destroy message
	DestroyMessage()

	//mark message as read
	MarkMessageAsRead()

	//mark message as received
	MarkMessageAsReceived()

	//mark message as destroyed
	MarkMessageAsDestroyed()

	//get messages for user
	GetMessages()
}
