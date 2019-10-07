package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/session"
)

// CreateCourt method
func CreateCourt(token string, c *court.Court) (string, error) {

	session := session.LookupToken(token)
	if session == nil {
		return "", codeerror.NewUnauthorized("Not Authorised")
	}

	// Check there are not too many players on the court
	info, err := court.GetCourtInfo()
	if err != nil {
		return "", err
	}

	if len(c.Container.Players) > info.PlayersPerCourt {
		return "", codeerror.NewBadRequest("Too many players on court")
	}

	// Check the people on the court are valid
	for _, id := range c.Container.Players {
		p, err := person.Load(id)
		if err != nil {
			return "", err
		}

		if p == nil {
			return "", codeerror.NewBadRequest(fmt.Sprintf("person[%s] not found", id))
		}

		if p.Player == false {
			return "", codeerror.NewBadRequest(fmt.Sprintf("person[%s] is not a player", id))
		}
	}

	// Save the court to disk
	id, err := c.Add()
	if err != nil {
		return "", err
	}

	return id, nil
}
