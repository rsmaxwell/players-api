package players

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/rsmaxwell/players-api/jsonTypes"
	"github.com/rsmaxwell/players-api/logger"
)

// RegisterRequest structure
type RegisterRequest struct {
	UserID    string `json:"userID"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Person Structure
type Person struct {
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	HashedPassword []byte `json:"hashedpassword"`
	Player         bool   `json:"player"`
}

// JSONPerson Structure
type JSONPerson struct {
	FirstName jsonTypes.JSONString `json:"firstname"`
	LastName  jsonTypes.JSONString `json:"lastname"`
	UserID    jsonTypes.JSONString `json:"userID"`
	Player    jsonTypes.JSONBool   `json:"player"`
}

var (
	peopleDirectory     string
	peopleDataDirectory string
)

// CreatePeopleDirectory  creates the people directory
func CreatePeopleDirectory() error {

	_, err := os.Stat(peopleDirectory)
	if err != nil {
		err := os.MkdirAll(peopleDirectory, 0755)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(peopleDataDirectory)
	if err != nil {
		err := os.MkdirAll(peopleDataDirectory, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemovePeopleDirectory removes ALL the person files
func RemovePeopleDirectory() error {
	logger.Logger.Printf("Remove person directory")

	_, err := os.Stat(peopleDirectory)
	if err == nil {
		err := removeContents(peopleDirectory)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(peopleDirectory)
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
	return person, nil
}

// UpdatePerson update fields
func UpdatePerson(id string, person2 JSONPerson) (*Person, error) {

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

	return person, nil
}

// RegisterPerson adds a person to the list of people
func RegisterPerson(reg RegisterRequest) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Printf("%s", err)
		return err
	}

	p, err := NewPerson(reg.FirstName, reg.LastName, hashedPassword, false)
	if err != nil {
		log.Println(err)
		return err
	}

	err = AddPerson(reg.UserID, *p)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// AddPerson adds a person to the list of people
func AddPerson(userID string, person Person) error {

	logger.Logger.Printf("AddPerson")

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	personJSON, err := json.Marshal(person)
	if err != nil {
		logger.Logger.Print(err)
		return err
	}

	logger.Logger.Printf("AddPerson: %s", personJSON)

	ok := checkID(userID)
	if !ok {
		return fmt.Errorf("The UserID [%s] is not valid", userID)
	}

	filename := peopleDataDirectory + "/" + userID + ".json"
	_, err = os.Stat(filename)
	if err == nil {
		return fmt.Errorf("UserId [%s] already exists", userID)
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

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	files, err := ioutil.ReadDir(peopleDataDirectory)
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

// GetPersonDetails returns the details of the person with the matching ID
func GetPersonDetails(id string) (*Person, error) {

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	personfile := peopleDataDirectory + "/" + id + ".json"
	if _, err := os.Stat(personfile); os.IsNotExist(err) {
		logger.Logger.Printf("The person file was not found. err = %s\n", err)
		return nil, err
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

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	personfile := peopleDataDirectory + "/" + id + ".json"
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

	err = CreatePeopleDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}
