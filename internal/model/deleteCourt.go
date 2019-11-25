package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionDeleteCourt = debug.NewFunction(pkg, "DeleteCourt")
)

// DeleteCourt method
func DeleteCourt(id string) error {
	f := functionDeleteCourt
	f.DebugVerbose("id:%s", id)

	err := court.Remove(id)
	if err != nil {
		return err
	}

	return nil
}
