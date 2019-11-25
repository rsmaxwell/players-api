package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionUpdateCourt = debug.NewFunction(pkg, "UpdateCourt")
)

// UpdateCourt method
func UpdateCourt(id string, fields map[string]interface{}) error {
	f := functionUpdateCourt
	f.DebugVerbose("id: %s, fields: %v", id, fields)

	ref := &common.Reference{Type: "court", ID: id}
	err := court.Update(ref, fields)
	if err != nil {
		return err
	}

	return nil
}
