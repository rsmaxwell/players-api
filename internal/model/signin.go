package model

// Registration type
type Signin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

// NewRegistration initialises a Signin object
func NewSignin(username string, password string) *Signin {
	r := new(Signin)
	r.Username = username
	r.Password = password
	return r
}
