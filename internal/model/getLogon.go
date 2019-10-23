package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/session"
)

// Login method
func Login(id, password string) (string, error) {

	p, err := person.Load(id)
	if err != nil {
		return "", codeerror.NewUnauthorized("Invalid userID and/or password")
	}

	if !p.CheckPassword(password) {
		return "", codeerror.NewUnauthorized("Invalid userID and/or password")
	}

	if !p.CanLogin() {
		return "", codeerror.NewUnauthorized("Not Authorized")
	}

	token, err := session.New(id)
	if err != nil {
		return "", err
	}

	return token, nil
}
