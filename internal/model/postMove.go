package model

import (
	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/destination"
	"github.com/rsmaxwell/players-api/internal/basic/queue"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionPostMove = debug.NewFunction(pkg, "PostMove")
)

// PostMove method
func PostMove(source, target *common.Reference, players []string) error {
	f := functionPostMove
	f.DebugVerbose("source: %v, target: %v, players: %v", source, target, players)

	// **********************************************************
	// * Load the source and target
	// **********************************************************
	s, err := load(source)
	if err != nil {
		return err
	}

	t, err := load(target)
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

// load method
func load(ref *common.Reference) (destination.Destination, error) {

	if ref.Type == "court" {
		return court.Load(ref)
	}

	return queue.Load(ref)
}
