package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionListPeople = debug.NewFunction(pkg, "ListPeople")
)

// ListPeople method
func ListPeople(filter []string) ([]string, error) {
	f := functionListPeople
	f.DebugVerbose("Filter: %s", filter)

	listOfPeople, err := person.List(filter)
	if err != nil {
		return nil, err
	}

	return listOfPeople, nil
}
