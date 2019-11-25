package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionAuthenticate = debug.NewFunction(pkg, "Authenticate")
)

// Authenticate method
func Authenticate(id, password string) error {
	f := functionAuthenticate
	f.DebugVerbose("id: %s, password:%s", id, "********")

	p, err := person.Load(id)
	if err != nil {
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CheckPassword(password) {
		f.Dump("password check failed")
		return codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CanLogin() {
		f.Dump("person not authorized to login")
		return codeerror.NewUnauthorized("Not Authorized")
	}

	return nil
}
