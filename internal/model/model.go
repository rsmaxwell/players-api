package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/common"
)

// Startup checks the state on disk is consistent
func Startup() error {

	// Make a list of players
	listOfPeople, err := person.List(person.AllRoles)
	if err != nil {
		return err
	}

	listOfPlayers := []string{}
	for _, id := range listOfPeople {
		p, err := person.Load(id)
		if err != nil {
			return err
		}
		if p.Player {
			listOfPlayers = append(listOfPlayers, id)
		}
	}

	// Subtract the players on courts away from the list of players
	listOfCourts, err := court.List()
	if err != nil {
		return err
	}
	for _, id := range listOfCourts {
		ref := common.Reference{Type: "court", ID: id}
		c, err := court.Load(&ref)
		if err != nil {
			return err
		}

		text := fmt.Sprintf("Court[%s]", id)
		listOfPlayers, err = common.SubtractLists(listOfPlayers, c.Container.Players, text)
		if err != nil {
			return err
		}
	}

	// Subtract the players waiting in the queue away from the list of players
	ref := common.Reference{Type: "queue", ID: ""}
	q, err := queue.Load(&ref)
	if err != nil {
		return err
	}

	listOfPlayers, err = common.SubtractLists(listOfPlayers, q.Container.Players, "Queue")
	if err != nil {
		return err
	}

	// The list of players should now be empty, however add any remaining players to the waiting queue
	for _, id := range listOfPlayers {
		q.Container.Players = append(q.Container.Players, id)
	}

	// Save the updated queue
	err = q.Save(&ref)
	if err != nil {
		return err
	}

	return nil
}
