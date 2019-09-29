package destination

import (
	"fmt"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
)

// Destination is the Generic Destination interface
type Destination interface {
	Load(ref *Reference) (*PeopleContainer, error)
	GetContainer() *PeopleContainer
	Show(title string)
	CheckPlayersLocation(players []string) error
	CheckSpace(players []string) error
	RemovePlayers(players []string) error
	AddPlayers(players []string) error
	Save(ref *Reference) error
	Update(fields map[string]interface{}) error
}

// Reference type
type Reference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

var (
	baseDir string
)

func init() {
	baseDir = common.RootDir
}

// FormatReference function
func FormatReference(ref *Reference) string {
	if ref.Type == "court" {
		return fmt.Sprintf("court[%s]", ref.ID)
	}

	return "queue"
}

// Load method
func Load(ref *Reference) (Destination, error) {

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
			return codeError.NewBadRequest(fmt.Sprintf("Player[%s] not found in source of Move command", personID))
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

// Show method
func Show(title string) {
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
func Save(ref *Reference) error {
	panic("This is an abstract method. Please provide an implementation!")
}

// Update updates the fields of the destination
func Update(fields map[string]interface{}) error {
	panic("This is an abstract method. Please provide an implementation!")
}
