package peoplecontainer

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
)

// PeopleContainer Structure
type PeopleContainer struct {
	Name    string   `json:"name" validate:"required,min=1,max=20"`
	Players []string `json:"players" validate:"dive,min=1,max=8"`
}

// NewContainer initialises a PeopleContainer object
func NewContainer(name string) *PeopleContainer {
	court := new(PeopleContainer)
	court.Name = name
	return court
}

// Update method
func (c *PeopleContainer) Update(container2 map[string]interface{}) error {

	if v, ok := container2["Name"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest(fmt.Sprintf("'Name' was an unexpected type. Expected: 'string'. Actual: %T", v))
		}
		c.Name = value
	}

	if v, ok := container2["Players"]; ok {

		var err error

		if value, ok := v.([]interface{}); ok {

			array := []string{}
			for _, i := range value {

				id, ok := i.(string)
				if !ok {
					return codeerror.NewBadRequest(fmt.Sprintf("'Player' was an unexpected type. Expected: 'string'. Actual: %T", v))
				}

				array, err = updateListOfPeople(id, array)
				if err != nil {
					return err
				}
			}
			c.Players = array

		} else if value, ok := v.([]string); ok {

			array := []string{}
			for _, id := range value {
				array, err = updateListOfPeople(id, array)
				if err != nil {
					return err
				}
			}
			c.Players = array

		} else {
			return codeerror.NewBadRequest(fmt.Sprintf("'Players' was an unexpected type. Expected: '[]interface{}'. Actual: %T %v", v, v))
		}
	}

	return nil
}

func updateListOfPeople(id string, array []string) ([]string, error) {

	p, err := person.Load(id)
	if err != nil {
		return nil, err
	}

	if !p.IsPlayer() {
		return array, codeerror.NewBadRequest(fmt.Sprintf("Person [%s] is not a player", id))
	}

	return append(array, id), nil
}

// EqualContainer returns 'true' if the containers are equal
func EqualContainer(c1, c2 PeopleContainer) bool {

	if c1.Name != c2.Name {
		return false
	}

	if !common.EqualArrayOfStrings2(c1.Players, c2.Players) {
		return false
	}

	return true
}
