package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionCreateCourt = debug.NewFunction(pkg, "CreateCourt")
)

// CreateCourt method
func CreateCourt(c *court.Court) (string, error) {
	f := functionCreateCourt
	f.DebugVerbose("court:%v", c)

	// Check there are not too many players on the court
	info, err := court.GetCourtInfo()
	if err != nil {
		return "", err
	}

	if len(c.Container.Players) > info.PlayersPerCourt {
		message := fmt.Sprintf("Too many players on court")
		f.DebugVerbose(message)
		return "", codeerror.NewBadRequest(message)
	}

	// Check the people on the court are valid
	for _, id := range c.Container.Players {
		p, err := person.Load(id)
		if err != nil {
			return "", err
		}

		if p == nil {
			message := fmt.Sprintf("person[%s] not found", id)
			f.DebugVerbose(message)
			return "", codeerror.NewBadRequest(message)
		}

		if p.Player == false {
			message := fmt.Sprintf("person[%s] is not a player", id)
			f.DebugVerbose(message)
			return "", codeerror.NewBadRequest(message)
		}
	}

	// Save the court to disk
	id, err := c.Add()
	if err != nil {
		return "", err
	}

	return id, nil
}
