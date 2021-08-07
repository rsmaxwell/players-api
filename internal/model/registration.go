package model

import (
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"

	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

// Registration type
type Registration struct {
	FirstName string `json:"firstname" validate:"required,min=3,max=20"`
	LastName  string `json:"lastname" validate:"required,min=3,max=20"`
	Knownas   string `json:"displayname" validate:"required,min=2,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"max=20"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
}

var (
	functionToPerson = debug.NewFunction(pkg, "ToPerson")
)

var (
	validate = validator.New()
)

// NewRegistration initialises a Registration object
func NewRegistration(firstname string, lastname string, knownas string, email string, phone string, password string) *Registration {
	r := new(Registration)
	r.FirstName = firstname
	r.LastName = lastname
	r.Knownas = knownas
	r.Email = email
	r.Phone = phone
	r.Password = password
	return r
}

// ToPerson converts a Registration into a person
func (r *Registration) ToPerson() (*FullPerson, error) {
	f := functionToPerson

	err := validate.Struct(r)
	if err != nil {
		message := fmt.Sprintf("validation failed for [%s]: %s", r.Email, err.Error())
		f.DebugVerbose(message)
		return nil, codeerror.NewBadRequest(err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.MinCost)
	if err != nil {
		message := "Could not generate password hash"
		f.Errorf(message)
		f.DumpError(err, message)
		return nil, err
	}
	p := NewPerson(r.FirstName, r.LastName, r.Knownas, r.Email, r.Phone, hash)

	return p, nil
}
