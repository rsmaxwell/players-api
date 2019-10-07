package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// GetCourt method
func GetCourt(token, id string) (*court.Court, error) {

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	ref := &common.Reference{Type: "court", ID: id}
	court, err := court.Load(ref)
	if err != nil {
		return nil, err
	}

	return court, nil
}
