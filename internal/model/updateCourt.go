package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// UpdateCourt method
func UpdateCourt(token string, id string, fields map[string]interface{}) error {

	session := session.LookupToken(token)
	if session == nil {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	if !person.CanUpdateCourt(session.UserID) {
		return codeerror.NewUnauthorized("Not Authorised")
	}

	ref := &common.Reference{Type: "court", ID: id}
	err := court.Update(ref, fields)
	if err != nil {
		return err
	}

	return nil
}