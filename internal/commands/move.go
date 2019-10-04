package commands

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

// Move method
func Move(source, target *common.Reference, players []string) error {

	fmt.Printf("x")

	// **********************************************************
	// * Load the source and target
	// **********************************************************
	s, err := model.Load(source)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	t, err := model.Load(target)
	if err != nil {
		return err
	}

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
func EqualsContainerReference(a, b *common.Reference) bool {

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
