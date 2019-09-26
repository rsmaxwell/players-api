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

// Person type
type Person struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashedpassword"`
	Status         string `json:"status"`
	Player         bool   `json:"player"`
}

var (
	baseDir string
	listDir string
)

func init() {
	baseDir = common.RootDir + "/people"
	listDir = baseDir + "/list"
}

// removeAll removes ALL the people
func removeAll() error {

	_, err := os.Stat(listDir)
	if err == nil {
		err = common.RemoveContents(listDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// Clear all the people
func Clear() error {

	err := removeAll()
	if err != nil {
		return err
	}

	err = createDirs()
	if err != nil {
		return err
	}

	return nil
}

// makeFilename function
func makeFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	filename := listDir + "/" + id + ".json"
	return filename, nil
}

// createDirs  creates the people directory
func createDirs() error {

	_, err := os.Stat(listDir)
	if err != nil {
		err := os.MkdirAll(listDir, 0755)
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

	// Save the updated person to disk
	err = Save(id, person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

// Insert adds a person to the list
func Insert(id string, person *Person) error {

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(listDir)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}
	if len(files) == 0 {
		person.Status = "admin"
	}

	// Check the user does not already exist
	found := Exists(id)
	if found {
		return codeError.NewInternalServerError(fmt.Sprintf("Person [%s] already exists", id))
	}

	// Save the updated court to disk
	err = Save(id, person)
	if err != nil {
		return err
	}

	return nil
}

// Save writes a Person to disk
func Save(id string, person *Person) error {

	personJSON, err := json.Marshal(person)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// List returns a list of the person IDs
func List() ([]string, error) {

	files, err := ioutil.ReadDir(listDir)
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

// Load Load returns the Person with the given ID
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

	files, err := ioutil.ReadDir(listDir)
	if err != nil {
		return 0, codeError.NewInternalServerError(err.Error())
	}

	return len(files), nil
}
