package destination

import (
	"fmt"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/person"
	"github.com/rsmaxwell/players-api/utilities"
)

// Container Structure
type Container struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}

// NewContainer initialises a Container object
func NewContainer(name string) *Container {
	court := new(Container)
	court.Name = name
	return court
}

// Update method
func (c Container) Update(container2 map[string]interface{}) error {

	if v, ok := container2["Name"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeError.NewBadRequest(fmt.Sprintf("'Name' was an unexpected type. Expected: 'string'. Actual: %T", v))
		}
		c.Name = value
	}

	if v, ok := container2["Players"]; ok {

		array := []string{}
		for _, i := range v.([]interface{}) {

			id2, ok := i.(string)
			if !ok {
				return codeError.NewBadRequest(fmt.Sprintf("'Player' was an unexpected type. Expected: 'string'. Actual: %T", i))
			}

			if !person.Exists(id2) {
				return codeError.NewNotFound(fmt.Sprintf("Person [%s] not found", id2))
			}

			if !person.IsPlayer(id2) {
				return codeError.NewBadRequest(fmt.Sprintf("Person [%s] is not a player", id2))
			}

			array = append(array, id2)
		}
		c.Players = array
	}

	return nil
}

// EqualContainer returns 'true' if the containers are equal
func EqualContainer(c1, c2 Container) bool {
	return utilities.Equal2(c1.Players, c2.Players)
}
