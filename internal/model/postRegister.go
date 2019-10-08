package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"golang.org/x/crypto/bcrypt"
)

// Register method
func Register(id, password, firstname, lastname, email string) error {

	if person.Exists(id) {
		return codeerror.NewBadRequest(fmt.Sprintf("Person[%s] already exists", id))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	list, err := person.List(person.AllRoles)
	if err != nil {
		return err
	}
	if len(list) >= 100 {
		return codeerror.NewBadRequest("Too many people")
	}

	err = person.New(firstname, lastname, email, hashedPassword, false).Save(id)
	if err != nil {
		return err
	}

	return nil
}
