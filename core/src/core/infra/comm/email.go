package comm

type EmailNotifier interface {
	//TODO: enviar uma mensagem com uma url de confirmação
	SendMessage()
}
