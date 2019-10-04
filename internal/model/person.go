package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// Person type
type Person struct {
	FirstName      string `json:"firstname" validate:"required,min=3,max=20"`
	LastName       string `json:"lastname" validate:"required,min=3,max=20"`
	Email          string `json:"email" validate:"required,email"`
	HashedPassword []byte `json:"hashedpassword" validate:"required,len=60"`
	Role           string `json:"role" validate:"required,oneof=admin normal suspended"`
	Player         bool   `json:"player"`
}

var (
	personBaseDir string
	personListDir string

	// RoleAdmin is allowed to do anything!
	RoleAdmin string

	// RoleNormal can change themselves
	RoleNormal string

	// RoleSuspended can do nothing. Only the 'admin' can change a suspended person
	RoleSuspended string

	// AllRoles is the 'all' filter which returns every person
	AllRoles []string
)

func init() {
	personBaseDir = common.RootDir + "/people"
	personListDir = personBaseDir + "/list"

	RoleAdmin = "admin"
	RoleNormal = "normal"
	RoleSuspended = "suspended"

	AllRoles = []string{RoleAdmin, RoleNormal, RoleSuspended}

	validate = validator.New()
}

// removeAll removes ALL the people
func removeAll() error {

	_, err := os.Stat(personListDir)
	if err == nil {
		err = common.RemoveContents(personListDir)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// makePersonFilename function
func makePersonFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	filename := personListDir + "/" + id + ".json"
	return filename, nil
}

// createDirs  creates the people directory
func createPersonDirs() error {

	_, err := os.Stat(personListDir)
	if err != nil {
		err := os.MkdirAll(personListDir, 0755)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// CheckUser - Basic check on the user calling the service
func CheckUser(id, password string) bool {

	p, err := LoadPerson(id)
	if err != nil {
		return false
	}
	if p == nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(p.HashedPassword, []byte(password))
	if err != nil {
		return false
	}

	return true
}

// NewPerson initialises a Person
func NewPerson(firstname string, lastname string, email string, hashedPassword []byte, player bool) *Person {
	person := new(Person)
	person.FirstName = firstname
	person.LastName = lastname
	person.Email = email
	person.HashedPassword = hashedPassword
	person.Player = player
	person.Role = RoleSuspended
	return person
}

// UpdatePerson method
func UpdatePerson(id string, session *session.Session, fields map[string]interface{}) error {

	p, err := LoadPerson(id)
	if err != nil {
		return err
	}

	err = p.Update(session, fields)
	if err != nil {
		return err
	}

	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// Update method
func (p *Person) Update(session *session.Session, person2 map[string]interface{}) error {

	if v, ok := person2["FirstName"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.FirstName' should be a string")
		}
		p.FirstName = value
	}

	if v, ok := person2["LastName"]; ok {
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.LastName' should be a string")
		}
		p.LastName = value
	}

	if v, ok := person2["Player"]; ok {
		value, ok := v.(bool)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.Player' should be a bool")
		}
		p.Player = value
	}

	if v, ok := person2["Status"]; ok {

		// Check we have the authority to perform this update
		myself, err := LoadPerson(session.UserID)
		if err != nil {
			return err
		}

		// Only 'admin' users can change the status of a user
		if myself.Role == RoleAdmin {
			return codeerror.NewUnauthorized("Not authorised to update 'Person.Status'")
		}

		// Check the type of the new status is a string
		value, ok := v.(string)
		if !ok {
			return codeerror.NewBadRequest("The type of 'Person.Status' should be a string")
		}

		p.Role = value
	}

	return nil
}

// Add adds a person to the list
func (p *Person) Add(id string) error {

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	if len(files) == 0 {
		p.Role = RoleAdmin
	}

	// Check the user does not already exist
	found := PersonExists(id)
	if found {
		return codeerror.NewInternalServerError(fmt.Sprintf("Person [%s] already exists", id))
	}

	// Save the updated court to disk
	err = p.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// Save writes a Person to disk
func (p *Person) Save(id string) error {

	err := validate.Struct(p)
	if err != nil {
		return codeerror.NewBadRequest(err.Error())
	}

	personJSON, err := json.Marshal(p)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	filename, err := makePersonFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	return nil
}

// ListPeople returns a list of the person IDs with one of the allowed status values
func ListPeople(filter []string) ([]string, error) {

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	list := []string{}
	for _, filenameInfo := range files {
		filename := filenameInfo.Name()
		id := strings.TrimSuffix(filename, path.Ext(filename))

		p, err := LoadPerson(id)
		if err != nil {
			return nil, err
		}
		if !common.Contains(filter, p.Role) {
			continue
		}

		list = append(list, id)
	}

	return list, nil
}

// PersonExists returns 'true' if the person exists
func PersonExists(id string) bool {

	filename, err := makePersonFilename(id)
	if err != nil {
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

// PersonIsPlayer returns 'true' if the person exists and is a player
func PersonIsPlayer(id string) bool {

	person, err := LoadPerson(id)
	if err != nil {
		return false
	}
	if person == nil {
		return false
	}

	return person.Player
}

// LoadPerson returns the Person with the given ID
func LoadPerson(id string) (*Person, error) {

	filename, err := makePersonFilename(id)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeerror.NewNotFound(err.Error())
		}
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}
	return &p, nil
}

// RemovePerson the person with the given ID
func RemovePerson(id string) error {

	filename, err := makePersonFilename(id)
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)

	if err == nil { // File exists
		err = os.Remove(filename)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
		return nil

	} else if os.IsNotExist(err) { // File does not exist
		return codeerror.NewNotFound(fmt.Sprintf("File Not Found: %s", filename))
	}

	return codeerror.NewInternalServerError(err.Error())
}

// PeopleSize returns the number of people
func PeopleSize() (int, error) {

	files, err := ioutil.ReadDir(personListDir)
	if err != nil {
		return 0, codeerror.NewInternalServerError(err.Error())
	}

	return len(files), nil
}
