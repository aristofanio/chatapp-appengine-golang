package auth

//Manager manages the authetication
type Manager interface {
	Auth(uEmail, uPass string) (string, error)
	DeAuth(sToken string) error
	CheckAuth(sToken string) (bool, error)
}
