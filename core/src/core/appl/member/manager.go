package member

type Manager interface {

	//create a member
	Create()

	//destroy a member
	Destroy()

	//update a member
	Update()

	//valid a member
	Valid()

	//accept an other member
	BlockOther()

	//list others
	ListOthers()
}
