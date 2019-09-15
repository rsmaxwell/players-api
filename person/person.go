package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/jsonTypes"
	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/session"
)

// Person Structure
type Person struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	HashedPassword []byte `json:"hashedpassword"`
	Status         string `json:"status"`
	Player         bool   `json:"player"`
}

// UpdatePersonRequest structure
type UpdatePersonRequest struct {
	Token  string     `json:"token"`
	Person JSONPerson `json:"person"`
}

// JSONPerson Structure
type JSONPerson struct {
	UserID    jsonTypes.JSONString `json:"userID"`
	FirstName jsonTypes.JSONString `json:"firstname"`
	LastName  jsonTypes.JSONString `json:"lastname"`
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
	logger.Logger.Printf("peopleDirectory = %s\n", peopleDir)
}

// CreatePeopleDirs  creates the people directory
func CreatePeopleDirs() error {

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

// RemovePeopleDirectory removes ALL the person files
func RemovePeopleDirectory() error {
	logger.Logger.Printf("Remove person directory")

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

// NewPerson initialises a Person object
func NewPerson(firstname string, lastname string, hashedPassword []byte, player bool) (*Person, error) {
	person := new(Person)
	person.FirstName = firstname
	person.LastName = lastname
	person.HashedPassword = hashedPassword
	person.Player = player
	person.Status = "normal"
	return person, nil
}

// UpdatePerson update fields
func UpdatePerson(id string, session *session.Session, person2 JSONPerson) (*Person, error) {

	logger.Logger.Printf("UpdatePerson")

	person, err := GetPersonDetails(id)
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
		myself, err := GetPersonDetails(session.UserID)
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

// AddPerson adds a person to the list of people
func AddPerson(userID string, person Person) error {

	filename := peopleListDir + "/" + userID + ".json"
	_, err := os.Stat(filename)
	if err == nil {
		return fmt.Errorf("UserId [%s] already exists", userID)
	}

	// Check the chearacters in the userid are sensible
	ok := common.CheckID(userID)
	if !ok {
		return fmt.Errorf("The UserID [%s] is not valid", userID)
	}

	// Make sure the people dirs exist
	err = CreatePeopleDirs()
	if err != nil {
		log.Println(err)
		return err
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

	err = ioutil.WriteFile(filename, personJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	return nil
}

// ListAllPeople returns a list of the person IDs
func ListAllPeople() ([]int, error) {

	err := CreatePeopleDirs()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	files, err := ioutil.ReadDir(peopleListDir)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	list := []int{}
	for _, filenameInfo := range files {
		var filename = filenameInfo.Name()
		var extension = filepath.Ext(filename)
		var name = filename[0 : len(filename)-len(extension)]

		id, err := strconv.Atoi(name)

		if err != nil {
			logger.Logger.Printf("Skipping unexpected person filename \"%s\". err = %s\n", filename, err)
		}

		list = append(list, id)
	}

	return list, nil
}

// PersonExists returns 'true' if the person exists
func PersonExists(id string) bool {

	personfile := peopleListDir + "/" + id + ".json"
	if _, err := os.Stat(personfile); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetPersonDetails returns the details of the person with the matching ID
func GetPersonDetails(id string) (*Person, error) {

	err := CreatePeopleDirs()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	personfile := peopleListDir + "/" + id + ".json"
	if _, err := os.Stat(personfile); os.IsNotExist(err) {
		logger.Logger.Printf("The person file was not found. err = %s\n", err)
		return nil, nil
	}

	data, err := ioutil.ReadFile(personfile)
	if err != nil {
		logger.Logger.Printf("Could not read file. err = %s\n", err)
		return nil, err
	}

	var p Person
	err = json.Unmarshal(data, &p)
	if err != nil {
		logger.Logger.Printf("Could not parse info data. err = %s\n", err)
		return nil, err
	}
	return &p, nil
}

// DeletePerson the person with the matching ID
func DeletePerson(id string) error {

	err := CreatePeopleDirs()
	if err != nil {
		log.Println(err)
		return err
	}

	personfile := peopleListDir + "/" + id + ".json"
	_, err = os.Stat(personfile)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("person [%d] not found", id)
	}

	err = os.Remove(personfile)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("could not delete person [%d]", id)
	}

	return nil
}

// ClearPeople resets the list of people
func ClearPeople() error {

	err := RemovePeopleDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = CreatePeopleDirs()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}
