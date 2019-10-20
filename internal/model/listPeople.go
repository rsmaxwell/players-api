package model

import (
	"log"

	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/session"
)

// ListPeople method
func ListPeople(token string, filter []string) ([]string, error) {

	log.Printf("model.ListPeople:")
	log.Printf("    Filter:    %s", filter)
	for _, v := range filter {
		log.Printf("        %s", v)
	}

	session := session.LookupToken(token)
	if session == nil {
		return nil, codeerror.NewUnauthorized("Not Authorised")
	}

	listOfPeople, err := person.List(filter)
	if err != nil {
		return nil, err
	}

	return listOfPeople, nil
}
