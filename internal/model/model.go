package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/common"

	"gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

// Startup checks the state on disk is consistent
func Startup() error {

	// Make a list of players
	listOfPeople, err := ListPeople(AllRoles)
	if err != nil {
		return err
	}

	listOfPlayers := []string{}
	for _, id := range listOfPeople {
		p, err := LoadPerson(id)
		if err != nil {
			return err
		}
		if p.Player {
			listOfPlayers = append(listOfPlayers, id)
		}
	}

	// Subtract the players on courts away from the list of players
	listOfCourts, err := ListCourts()
	if err != nil {
		return err
	}
	for _, id := range listOfCourts {
		ref := common.Reference{Type: "court", ID: id}
		c, err := LoadCourt(&ref)
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
	q, err := LoadQueue(&ref)
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
