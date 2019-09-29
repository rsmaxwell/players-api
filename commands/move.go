package commands

import (
	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/destination"
)

// Move method
func Move(source, target *destination.Reference, players []string) error {

	// **********************************************************
	// * Load the source and target
	// **********************************************************
	s, err := destination.Load(source)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	t, err := destination.Load(target)
	if err != nil {
		return err
	}

	s.Show("source")
	t.Show("target")

	// **********************************************************
	// * Checks
	// **********************************************************
	err = s.CheckPlayersLocation(players)
	if err != nil {
		return err
	}

	err = t.CheckSpace(players)
	if err != nil {
		return err
	}

	// **********************************************************
	// * Move the players from source to target
	// **********************************************************
	err = s.RemovePlayers(players)
	if err != nil {
		return err
	}

	err = t.AddPlayers(players)
	if err != nil {
		return err
	}

	// **********************************************************
	// * Save the update source and targets to disk
	// **********************************************************
	err = s.Save(source)
	if err != nil {
		return err
	}

	err = t.Save(target)
	if err != nil {
		return err
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
