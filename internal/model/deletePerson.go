package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionDeletePerson = debug.NewFunction(pkg, "DeletePerson")
)

// DeletePerson method
func DeletePerson(id string) error {
	f := functionDeletePerson
	f.DebugVerbose("id:%s", id)

	err := person.Remove(id)
	if err != nil {
		return err
	}

	return nil
}
