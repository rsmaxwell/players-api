package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/session"
)

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(token, id string) (*person.Person, error) {
	f := functionGetPerson
	f.DebugVerbose("token: %s, id:%s", token, id)

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	p, err := person.Load(id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		common.MetricsData.ClientError++
		return nil, codeerror.NewNotFound(fmt.Sprintf("Person[%s] Not Found", id))
	}

	return p, nil
}
