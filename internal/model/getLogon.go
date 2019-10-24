package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionLogin = debug.NewFunction(pkg, "Login")
)

// Login method
func Login(id, password string) (string, error) {
	f := functionLogin
	f.DebugVerbose("id: %s, password:%s", id, "********")

	p, err := person.Load(id)
	if err != nil {
		return "", codeerror.NewUnauthorized("Invalid userID and/or password")
	}

	if !p.CheckPassword(password) {
		f.Dump("password check failed")
		return "", codeerror.NewUnauthorized("Invalid userID and/or password")
	}

	if !p.CanLogin() {
		f.Dump("person not authorized to login")
		return "", codeerror.NewUnauthorized("Not Authorized")
	}

	token, err := session.New(id)
	if err != nil {
		return "", err
	}

	return token, nil
}
