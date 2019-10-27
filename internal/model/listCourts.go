package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionListCourts = debug.NewFunction(pkg, "ListCourts")
)

// ListCourts method
func ListCourts() ([]string, error) {
	f := functionListCourts
	f.DebugVerbose("")

	listOfCourts, err := court.List()
	if err != nil {
		return nil, err
	}

	return listOfCourts, nil
}
