package invite

type Manager interface {

	//create an invite
	Create()

	//destroy an invite
	Destroy()

	//accept an invite
	Accept()
}
