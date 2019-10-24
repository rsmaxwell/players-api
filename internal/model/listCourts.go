package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionListCourts = debug.NewFunction(pkg, "ListCourts")
)

// ListCourts method
func ListCourts(token string) ([]string, error) {
	f := functionListCourts
	f.DebugVerbose("token: %s", token)

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	listOfCourts, err := court.List()
	if err != nil {
		return nil, err
	}

	return listOfCourts, nil
}
