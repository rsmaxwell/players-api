package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionGetPerson = debug.NewFunction(pkg, "GetPerson")
)

// GetPerson method
func GetPerson(id string) (*person.Person, error) {
	f := functionGetPerson
	f.DebugVerbose("id: %s", id)

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
