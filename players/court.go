package players

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rsmaxwell/players-api/logger"
)

// Court Structure
type Court struct {
	Name string `json:"name"`
}

var (
	courtDirectory     string
	courtDataDirectory string
	courtInfoFile      string
)

// CreateCourtDirectory  creates the people directory
func CreateCourtDirectory() error {

	_, err := os.Stat(courtDirectory)
	if err != nil {
		err := os.MkdirAll(courtDirectory, 0755)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(courtDataDirectory)
	if err != nil {
		err := os.MkdirAll(courtDataDirectory, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveCourtDirectory removes ALL the court files
func RemoveCourtDirectory() error {
	logger.Logger.Printf("Remove court directory")

	_, err := os.Stat(courtDirectory)
	if err == nil {
		err := removeContents(courtDirectory)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(courtDirectory)
	}

	return nil
}

// CreateCourtInfoFile initialises the Info object from file
func CreateCourtInfoFile() (*Info, error) {

	logger.Logger.Printf("Checking the infofile exists")
	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		err := CreateCourtDirectory()
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}

		info := Info{CurrentID: 1000}
		infoJSON, err := json.Marshal(info)
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}

		err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}
	}

	data, err := ioutil.ReadFile(courtInfoFile)
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

// GetAndIncrementCurrentCourtID returns the CurrentID and then increments the CurrentID on disk
func GetAndIncrementCurrentCourtID() (int, error) {

	// Make sure the court directory and info file exists
	_, err := CreateCourtInfoFile()
	if err != nil {
		return 0, err
	}

	// Read the court info file
	data, err := ioutil.ReadFile(courtInfoFile)
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

	err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
	if err != nil {
		return 0, err
	}

	return currentID, nil
}

// NewCourt initialises a Court object
func NewCourt(name string) (*Court, error) {
	court := new(Court)
	court.Name = name
	return court, nil
}

// AddCourt adds a court to the list of courts
func AddCourt(court Court) error {

	err := CreateCourtDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	id, err := GetAndIncrementCurrentCourtID()
	if err != nil {
		log.Println(err)
		return err
	}

	courtJSON, err := json.Marshal(court)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	err = CreateCourtDirectory()
	if err != nil {
		logger.Logger.Panicf(err.Error())
	}

	courtfile := courtDataDirectory + "/" + strconv.Itoa(id) + ".json"
	err = ioutil.WriteFile(courtfile, courtJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	return nil
}

// ListAllCourts returns a list of the court IDs
func ListAllCourts() ([]int, error) {

	err := CreateCourtDirectory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	files, err := ioutil.ReadDir(courtDataDirectory)
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
			logger.Logger.Printf("Skipping unexpected court filename \"%s\". err = %s\n", filename, err)
		}

		list = append(list, id)
	}

	return list, nil
}

// GetCourtDetails returns the details of the court with the matching ID
func GetCourtDetails(id int) (*Court, error) {

	err := CreateCourtDirectory()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	courtfile := courtDataDirectory + "/" + strconv.Itoa(id) + ".json"
	if _, err := os.Stat(courtfile); os.IsNotExist(err) {
		logger.Logger.Printf("The court file was not found. err = %s\n", err)
		return nil, err
	}

	data, err := ioutil.ReadFile(courtfile)
	if err != nil {
		logger.Logger.Printf("Could not read file. err = %s\n", err)
		return nil, err
	}

	var c Court
	err = json.Unmarshal(data, &c)
	if err != nil {
		logger.Logger.Printf("Could not parse info data. err = %s\n", err)
		return nil, err
	}
	return &c, nil
}

// DeleteCourt the court with the matching ID
func DeleteCourt(id int) error {

	err := CreateCourtDirectory()
	if err != nil {
		log.Println(err)
		return err
	}

	courtfile := courtDataDirectory + "/" + strconv.Itoa(id) + ".json"
	_, err = os.Stat(courtfile)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("court [%d] not found", id)
	}

	err = os.Remove(courtfile)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("could not delete court [%d]", id)
	}

	return nil
}

// ResetCourts resets the list of courts
func ResetCourts(list ...Court) error {

	err := RemoveCourtDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = CreateCourtDirectory()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	_, err = CreateCourtInfoFile()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	for _, court := range list {
		if err != nil {
			logger.Logger.Fatal(err)
		}

		err = AddCourt(court)
		if err != nil {
			logger.Logger.Fatal(err)
		}
	}

	return nil
}
