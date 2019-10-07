package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/session"
)

// DeletePerson method
func DeletePerson(token, id string) error {

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	err := person.Remove(id)
	if err != nil {
		return err
	}

	return nil
}
