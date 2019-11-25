package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionGetCourt = debug.NewFunction(pkg, "GetCourt")
)

// GetCourt method
func GetCourt(id string) (*court.Court, error) {
	f := functionGetCourt
	f.DebugVerbose("id:%s", id)

	ref := &common.Reference{Type: "court", ID: id}
	court, err := court.Load(ref)
	if err != nil {
		f.Dump("could not load court: %v", err)
		return nil, err
	}

	return court, nil
}
