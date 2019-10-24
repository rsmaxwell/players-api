package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionListPeople = debug.NewFunction(pkg, "ListPeople")
)

// ListPeople method
func ListPeople(token string, filter []string) ([]string, error) {
	f := functionListPeople
	f.DebugVerbose("Filter: %s", filter)

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	listOfPeople, err := person.List(filter)
	if err != nil {
		return nil, err
	}

	return listOfPeople, nil
}
