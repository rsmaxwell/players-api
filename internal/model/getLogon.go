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
func Authenticate(id, password string) (*person.Person, error) {
	f := functionAuthenticate
	f.DebugVerbose("id: %s, password:%s", id, "********")

	p, err := person.Load(id)
	if err != nil {
		f.DebugVerbose("could not load person [%s]", id)
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CheckPassword(password) {
		f.DebugVerbose("password check failed for person [%s]", id)
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	if !p.CanLogin() {
		f.DebugVerbose("person [%s] not authorized to login", id)
		return nil, codeerror.NewForbidden("Forbidden")
	}

	return p, nil
}
