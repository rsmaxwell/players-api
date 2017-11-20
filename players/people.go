package players

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/rsmaxwell/players-server/logger"
)

// Person Structure
type Person struct {
	Name string `json:"name"`
}

// Info structure
type Info struct {
	CurrentID int `json:"currentID"`
}

var (
	rootdir             string
	peopleDirectory     string
	peopleDataDirectory string
	peopleInfoFile      string
)

func homeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func init() {

	home := homeDir()

	if flag.Lookup("test.v") == nil {
		rootdir = home + "/players"
	} else {
		rootdir = home + "/players-test"
	}

	peopleDirectory = rootdir + "/people"
	peopleDataDirectory = peopleDirectory + "/data"
	peopleInfoFile = peopleDirectory + "/info.json"

	logger.Logger.Printf("peopleDirectory = %s\n", peopleDirectory)
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

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

// CreatePeopleInfoFile initialises the Info object from file
func CreatePeopleInfoFile() (*Info, error) {

	logger.Logger.Printf("Checking the infofile exists")
	if _, err := os.Stat(peopleInfoFile); os.IsNotExist(err) {

		err := CreatePeopleDirectory()
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}

		info := Info{CurrentID: 1000}
		infoJSON, err := json.Marshal(info)
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}

		err = ioutil.WriteFile(peopleInfoFile, infoJSON, 0644)
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}
	}

	data, err := ioutil.ReadFile(peopleInfoFile)
	if err != nil {
		logger.Logger.Panicf(err.Error())
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		logger.Logger.Panicf(err.Error())
	}

	return &i, nil
}

// GetAndIncrementCurrentID returns the CurrentID and then increments the CurrentID on disk
func GetAndIncrementCurrentID() (int, error) {

	// Make sure the person directory and info file exists
	_, err := CreatePeopleInfoFile()
	if err != nil {
		return 0, err
	}

	// Read the person info file
	data, err := ioutil.ReadFile(peopleInfoFile)
	if err != nil {
		return 0, err
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return 0, err
	}

	currentID := i.CurrentID

	i.CurrentID = i.CurrentID + 1

	infoJSON, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile(peopleInfoFile, infoJSON, 0644)
	if err != nil {
		return 0, err
	}

	return currentID, nil
}

// NewPerson initialises a Person object
func NewPerson(name string) (*Person, error) {
	person := new(Person)
	person.Name = name
	return person, nil
}

// AddPerson adds a person to the list of people
func AddPerson(person Person) error {

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := GetAndIncrementCurrentID()
	if err != nil {
		log.Println(err)
		return err
	}

	personJSON, err := json.Marshal(person)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	err = CreatePeopleDirectory()
	if err != nil {
		logger.Logger.Panicf(err.Error())
	}

	personfile := peopleDataDirectory + "/" + strconv.Itoa(id) + ".json"
	err = ioutil.WriteFile(personfile, personJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	return nil
}

// List returns a list of the person IDs
func List() ([]int, error) {

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

	var list []int
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

// Details returns the details of the person with the matching ID
func Details(id int) (*Person, error) {

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	personfile := peopleDataDirectory + "/" + strconv.Itoa(id) + ".json"
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

// Delete the person with the matching ID
func Delete(id int) error {

	err := CreatePeopleDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	personfile := peopleDataDirectory + "/" + strconv.Itoa(id) + ".json"
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

// Reset resets the list of people
func Reset(names ...string) error {

	err := RemovePeopleDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = CreatePeopleDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	_, err = CreatePeopleInfoFile()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	for _, element := range names {
		person, err := NewPerson(element)
		if err != nil {
			logger.Logger.Fatal(err)
		}

		err = AddPerson(*person)
		if err != nil {
			logger.Logger.Fatal(err)
		}
	}

	return nil
}
