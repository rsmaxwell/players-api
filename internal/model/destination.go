package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// Destination is the Generic Destination interface
type Destination interface {
	Load(ref *common.Reference) (*PeopleContainer, error)
	GetContainer() *PeopleContainer
	CheckPlayersLocation(players []string) error
	CheckSpace(players []string) error
	RemovePlayers(players []string) error
	AddPlayers(players []string) error
	Save(ref *common.Reference) error
	Update(session *session.Session, fields map[string]interface{}) error
}

// FormatReference function
func FormatReference(ref *common.Reference) string {
	if ref.Type == "court" {
		return fmt.Sprintf("court[%s]", ref.ID)
	}

	return "queue"
}

// Load method
func Load(ref *common.Reference) (Destination, error) {

	if ref.Type == "court" {
		return LoadCourt(ref)
	}

	return LoadQueue(ref)
}

// CheckPlayersInContainer checks the players are at this destination
func CheckPlayersInContainer(c *PeopleContainer, players []string) error {

	for _, personID := range players {
		found := false
		for _, id := range c.Players {
			if id == personID {
				found = true
				break
			}
		}
		if !found {
			return codeerror.NewBadRequest(fmt.Sprintf("Player[%s] not found in source of Move command", personID))
		}
	}
	return nil
}

// RemovePlayersFromContainer deletes players from the container
func RemovePlayersFromContainer(c *PeopleContainer, players []string) error {
	array := []string{}
	for _, id := range c.Players {

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

	c.Players = array
	return nil
}

// AddPlayersToContainer adds players to a container
func AddPlayersToContainer(c *PeopleContainer, players []string) error {
	array := []string{}
	for _, id := range c.Players {
		array = append(array, id)
	}

	for _, personID := range players {

		found := false
		for _, id := range c.Players {
			if id == personID {
				found = true
				break
			}
		}

		if !found {
			array = append(array, personID)
		}
	}

	c.Players = array
	return nil
}

// GetContainer method
func GetContainer() (*PeopleContainer, error) {
	panic("This is an abstract method. Please provide an implementation!")
}

// CheckPlayersLocation checks the players are at this destination
func CheckPlayersLocation(players []string) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// CheckSpace checks there is space in the target for the moving players
func CheckSpace(players []string) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// RemovePlayers deletes players from the destination
func RemovePlayers(players []string) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// AddPlayers adds players to the destination
func AddPlayers(players []string) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// Save writes the destination to disk
func Save(ref *common.Reference) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// Update updates the fields of the destination
func Update(fields map[string]interface{}) error {
	panic("This is an abstract method. Please provide an implementation!")
}
