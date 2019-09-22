package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/session"
)

// Person Structure
type Person struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashedpassword"`
	Status         string `json:"status"`
	Player         bool   `json:"player"`
}

var (
	peopleDir     string
	peopleListDir string
)

func init() {
	peopleDir = common.RootDir + "/people"
	peopleListDir = peopleDir + "/list"
}

// removeAllPeople removes ALL the person files
func removeAllPeople() error {

	_, err := os.Stat(peopleListDir)
	if err == nil {
		err = common.RemoveContents(peopleListDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// Clear - Just the list of people
func Clear() error {

	err := removeAllPeople()
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	err = createDirs()
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// makeFilename function
func makeFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	filename := peopleListDir + "/" + id + ".json"
	return filename, nil
}

// createDirs  creates the people directory
func createDirs() error {

	_, err := os.Stat(peopleListDir)
	if err != nil {
		err := os.MkdirAll(peopleListDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// New initialises a Person
func New(firstname string, lastname string, email string, hashedPassword []byte, player bool) *Person {
	person := new(Person)
	person.FirstName = firstname
	person.LastName = lastname
	person.Email = email
	person.HashedPassword = hashedPassword
	person.Player = player
	person.Status = "normal"
	return person
}

// Update method
func Update(id string, session *session.Session, person2 map[string]interface{}) (*Person, error) {

	var err error

	person, err := Load(id)
	if err != nil {
		return nil, err
	}

	if v, ok := person2["FirstName"]; ok {
		value, ok := v.(string)
		if !ok {
			return nil, codeError.NewBadRequest("The type of 'Person.FirstName' should be a string")
		}
		person.FirstName = value
	}

	if v, ok := person2["LastName"]; ok {
		value, ok := v.(string)
		if !ok {
			return nil, codeError.NewBadRequest("The type of 'Person.LastName' should be a string")
		}
		person.LastName = value
	}

	if v, ok := person2["Player"]; ok {
		value, ok := v.(bool)
		if !ok {
			return nil, codeError.NewBadRequest("The type of 'Person.Player' should be a bool")
		}
		person.Player = value
	}

	if v, ok := person2["Status"]; ok {

		// Check we have the authority to perform this update
		myself, err := Load(session.UserID)
		if err != nil {
			return nil, err
		}

		// Only 'admin' users can update the 'Status' field
		if myself.Status != "admin" {
			return nil, codeError.NewUnauthorized("Not authorised to update 'Person.Status'")
		}

		value, ok := v.(string)
		if !ok {
			return nil, codeError.NewBadRequest("The type of 'Person.Status' should be a string")
		}
		person.Status = value
	}

	// Convert the 'person' object into a JSON string
	personJSON, err := json.Marshal(person)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	// Save the updated person to disk
	filename, err := makeFilename(id)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	return person, nil
}

// Save adds a person to the list of people
func Save(id string, person *Person) error {

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}
	if len(files) == 0 {
		person.Status = "admin"
	}

	// Convert the 'person' object into a JSON string
	personJSON, err := json.Marshal(person)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	// Save the updated court to disk
	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	// Check the user does not already exist
	found := Exists(id)
	if found {
		return codeError.NewInternalServerError(fmt.Sprintf("Person [%s] already exists", filename))
	}

	// Save the person to disk
	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// List returns a list of the person IDs
func List() ([]string, error) {

	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	list := []string{}
	for _, filenameInfo := range files {
		filename := filenameInfo.Name()
		id := strings.TrimSuffix(filename, path.Ext(filename))
		list = append(list, id)
	}

	return list, nil
}

// Exists returns 'true' if the person exists
func Exists(id string) bool {

	filename, err := makeFilename(id)
	if err != nil {
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

// IsPlayer returns 'true' if the person exists and is a player
func IsPlayer(id string) bool {

	person, err := Load(id)
	if err != nil {
		return false
	}
	if person == nil {
		return false
	}

	return person.Player
}

// Load returns the details of the person with the given ID
func Load(id string) (*Person, error) {

	filename, err := makeFilename(id)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeError.NewNotFound(err.Error())
		}
		return nil, codeError.NewInternalServerError(err.Error())
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}
	return &p, nil
}

// Remove the person with the given ID
func Remove(id string) error {

	var err error

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}
	_, err = os.Stat(filename)

	if err == nil { // File exists
		err = os.Remove(filename)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
		return nil

	} else if os.IsNotExist(err) { // File does not exist
		return codeError.NewNotFound(fmt.Sprintf("File Not Found: %s", filename))
	}

	return codeError.NewInternalServerError(err.Error())
}

// Size returns the number of people
func Size() (int, error) {

	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		return 0, codeError.NewInternalServerError(err.Error())
	}

	return len(files), nil
}
