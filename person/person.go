package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/jsonTypes"
	"github.com/rsmaxwell/players-api/logger"
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

// JSONPerson Structure
type JSONPerson struct {
	UserID    jsonTypes.JSONString `json:"userID"`
	FirstName jsonTypes.JSONString `json:"firstname"`
	LastName  jsonTypes.JSONString `json:"lastname"`
	Email     jsonTypes.JSONString `json:"email"`
	Password  jsonTypes.JSONString `json:"password"`
	Status    jsonTypes.JSONString `json:"status"`
	Player    jsonTypes.JSONBool   `json:"player"`
}

var (
	peopleDir     string
	peopleListDir string
)

func init() {

	peopleDir = common.RootDir + "/people"
	peopleListDir = peopleDir + "/list"
}

// createDirs  creates the people directory
func createDirs() error {

	_, err := os.Stat(peopleDir)
	if err != nil {
		err := os.MkdirAll(peopleDir, 0755)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(peopleListDir)
	if err != nil {
		err := os.MkdirAll(peopleListDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// removeDir removes ALL the person files
func removeAllDirs() error {

	_, err := os.Stat(peopleDir)
	if err == nil {
		err := common.RemoveContents(peopleDir)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(peopleDir)
	}

	return nil
}

// removeDir removes ALL the person files
func removeListDir() error {

	_, err := os.Stat(peopleListDir)
	if err == nil {
		err := common.RemoveContents(peopleListDir)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(peopleDir)
	}

	return nil
}

// ClearAll - clears ALL the people directories
func ClearAll() error {

	err := removeAllDirs()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = createDirs()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}

// Clear - Just the list of people
func Clear() error {

	err := removeListDir()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = createDirs()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}

// New initialises a Person
func New(firstname string, lastname string, email string, hashedPassword []byte, player bool) (*Person, error) {
	person := new(Person)
	person.FirstName = firstname
	person.LastName = lastname
	person.Email = email
	person.HashedPassword = hashedPassword
	person.Player = player
	person.Status = "normal"
	return person, nil
}

// Update method
func Update(id string, session *session.Session, person2 JSONPerson) (*Person, error) {

	person, err := Get(id)
	if err != nil {
		logger.Logger.Print(err)
		return nil, fmt.Errorf("person [%d] not found", id)
	}

	if person2.FirstName.Set {
		person.FirstName = person2.FirstName.Value
	}

	if person2.LastName.Set {
		person.LastName = person2.LastName.Value
	}

	if person2.Player.Set {
		person.Player = person2.Player.Value
	}

	if person2.Status.Set {
		// Check we have the authority to perform this update
		myself, err := Get(session.UserID)
		if err != nil {
			logger.Logger.Print(err)
			return nil, fmt.Errorf("person [%d] not found", id)
		}

		// Only 'admin' users can update the 'Status' field
		if myself.Status != "admin" {
			logger.Logger.Print(err)
			return nil, fmt.Errorf("Not authorised")
		}
		person.Status = person2.Status.Value
	}

	// Convert the 'person' object into a JSON string
	personJSON, err := json.Marshal(person)
	if err != nil {
		logger.Logger.Print(err)
		return nil, err
	}

	// Save the updated person to disk
	filename := peopleListDir + "/" + id + ".json"
	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return nil, fmt.Errorf("internal error")
	}

	return person, nil
}

// Add adds a person to the list of people
func Add(id string, person Person) error {

	// Check the characters in the id are sensible
	ok := common.CheckID(id)
	if !ok {
		return fmt.Errorf("The id [%s] is not valid", id)
	}

	// Check the person does not already exist
	if Exists(id) {
		return fmt.Errorf("Person[%s] already exists", id)
	}

	// The first user must be made an 'admin' user
	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		logger.Logger.Print(err)
		return err
	}
	if len(files) == 0 {
		person.Status = "admin"
	}

	// Convert the 'person' object into a JSON string
	personJSON, err := json.Marshal(person)
	if err != nil {
		logger.Logger.Print(err)
		return err
	}

	filename := peopleListDir + "/" + id + ".json"
	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	return nil
}

// List returns a list of the person IDs
func List() ([]string, error) {

	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		log.Println(err)
		return nil, err
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

	filename := peopleListDir + "/" + id + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// Get returns the details of the person with the given ID
func Get(id string) (*Person, error) {

	filename := peopleListDir + "/" + id + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		logger.Logger.Printf("File not found. %s", filename)
		return nil, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Logger.Printf("Could not read file. err = %s", err)
		return nil, err
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		logger.Logger.Printf("Could not parse file. err = %s", err)
		return nil, err
	}
	return &p, nil
}

// Delete the person with the given ID
func Delete(id string) error {

	filename := peopleListDir + "/" + id + ".json"
	_, err := os.Stat(filename)

	if err == nil { // File exists
		err = os.Remove(filename)
		if err != nil {
			logger.Logger.Print(err)
			return err
		}
	} else if os.IsNotExist(err) { // File does not exist
		return nil
	} else {
		logger.Logger.Print(err)
		return err
	}

	return nil
}
