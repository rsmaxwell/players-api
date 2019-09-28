package commands

import (
	"fmt"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/destination"
)

// Move method
func Move(source, target *destination.Reference, players []string) error {

	var q1, q2 *destination.Queue
	var c1, c2 *destination.Court
	var sc1, sc2 *destination.Container
	var err error

	// **********************************************************
	// * De-reference the source and target
	// **********************************************************

	// De-reference the source
	if source.Type == "court" {
		c1, err = destination.LoadCourt(source.ID)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		sc1 = &c1.Container
	} else {
		q1, err = destination.LoadQueue()
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		sc1 = &q1.Container
	}

	// De-reference the target
	if target.Type == "court" {
		c2, err = destination.LoadCourt(target.ID)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		sc2 = &c2.Container
	} else {
		q2, err = destination.LoadQueue()
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		sc2 = &q2.Container
	}

	// **********************************************************
	// * Checks
	// **********************************************************

	// Check all the moving players are at the source
	for _, personID := range players {
		found := false
		for _, id := range sc1.Players {
			if id == personID {
				found = true
				break
			}
		}
		if !found {
			return codeError.NewBadRequest(fmt.Sprintf("Player[%s] not found in source of Move command", personID))
		}
	}

	// Check there is space in the target for the moving players
	if target.Type == "court" {
		info, err := destination.GetCourtInfo()
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		containerSize := len(c2.Container.Players)
		playersSize := len(players)

		if containerSize+playersSize > info.PlayersPerCourt {
			return codeError.NewBadRequest(fmt.Sprintf("Too many players. %d + %d > %d", containerSize, playersSize, info.PlayersPerCourt))
		}
	}

	// **********************************************************
	// * Delete the moving players from the source
	// **********************************************************
	array := []string{}
	for _, id := range sc1.Players {

		found := false
		for _, personID := range players {
			if id == personID {
				found = true
				break
			}
		}

		if !found {
			array = append(array, id)
		}
	}

	sc1.Players = array

	// **********************************************************
	// * Add the moving players to the target
	// **********************************************************

	array = []string{}
	for _, id := range sc2.Players {
		array = append(array, id)
	}

	for _, personID := range players {

		found := false
		for _, id := range sc2.Players {
			if id == personID {
				found = true
				break
			}
		}

		if !found {
			array = append(array, personID)
		}
	}

	sc2.Players = array

	// **********************************************************
	// * Save the update source and target to disk
	// **********************************************************

	// Save the updated soure
	if source.Type == "court" {
		err := c1.Save(source.ID)
		if err != nil {
			return err
		}
	} else {
		err := q1.Save()
		if err != nil {
			return err
		}
	}

	// Save the updated target
	if target.Type == "court" {
		err := c2.Save(target.ID)
		if err != nil {
			return err
		}
	} else {
		err := q2.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// EqualsContainerReference function
func EqualsContainerReference(a, b *destination.Reference) bool {

	if a.Type != b.Type {
		return false
	}

	if a.Type == "court" {
		if a.ID != b.ID {
			return false
		}
	}

	return true
}
