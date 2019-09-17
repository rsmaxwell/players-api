package court

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/jsonTypes"

	"github.com/rsmaxwell/players-api/logger"
)

// Info structure
type Info struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// Court Structure
type Court struct {
	Name    string   `json:"name"`
	Players []string `json:"players"`
}

// JSONCourt Structure
type JSONCourt struct {
	Name jsonTypes.JSONString `json:"name"`
}

var (
	courtDir      string
	courtListDir  string
	courtInfoFile string
)

func init() {

	courtDir = common.RootDir + "/courts"
	courtListDir = courtDir + "/list"
	courtInfoFile = courtDir + "/info.json"
}

// createDirs creates the people directory
func createDirs() error {

	_, err := os.Stat(courtDir)
	if err != nil {
		err := os.MkdirAll(courtDir, 0755)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(courtListDir)
	if err != nil {
		err := os.MkdirAll(courtListDir, 0755)
		if err != nil {
			return err
		}
	}

	_, err = createInfo()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}

// removeAllDirs - Remove ALL the court directories
func removeAllDirs() error {

	_, err := os.Stat(courtDir)
	if err == nil {
		err := common.RemoveContents(courtDir)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(courtDir)
	}

	return nil
}

// removeListDir - remove the list of courts
func removeListDir() error {

	_, err := os.Stat(courtListDir)
	if err == nil {
		err := common.RemoveContents(courtListDir)
		if err != nil {
			logger.Logger.Panic(err.Error())
		}

		os.Remove(courtDir)
	}

	return nil
}

// ClearAll - Clear ALL the court directories
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

// Clear - Clear the list of courts
func Clear() error {

	err := removeListDir()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = os.Remove(courtInfoFile)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	err = createDirs()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return nil
}

// New initialises a Court object
func New(name string) *Court {
	court := new(Court)
	court.Name = name
	return court
}

// Update method
func Update(id string, court2 JSONCourt) (*Court, error) {

	court, err := Get(id)
	if err != nil {
		logger.Logger.Print(err)
		return nil, fmt.Errorf("court [%s] not found", id)
	}

	if court2.Name.Set {
		court.Name = court2.Name.Value
	}

	return court, nil
}

// Add adds a court to the list of courts
func Add(court Court) error {

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		log.Println(err)
		return err
	}

	courtJSON, err := json.Marshal(court)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	id := strconv.Itoa(count)
	filename := courtListDir + "/" + id + ".json"
	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		logger.Logger.Print(err)
		return fmt.Errorf("internal error")
	}

	return nil
}

// List returns a list of the court IDs
func List() ([]string, error) {

	files, err := ioutil.ReadDir(courtListDir)
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

// Exists returns 'true' if the court exists
func Exists(id string) bool {

	filename := courtListDir + "/" + id + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// Get returns the details of the court with the given ID
func Get(id string) (*Court, error) {

	filename := courtListDir + "/" + id + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		logger.Logger.Printf("File not found. %s", filename)
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Logger.Printf("Could not read file. err = %s", err)
		return nil, err
	}

	var c Court
	err = json.Unmarshal(data, &c)
	if err != nil {
		logger.Logger.Printf("Could not parse file. err = %s", err)
		return nil, err
	}
	return &c, nil
}

// Delete the court with the given ID
func Delete(id string) error {

	filename := courtListDir + "/" + id + ".json"
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

// createInfo initialises the Info objects
func createInfo() (*Info, error) {

	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		info := new(Info)
		info.NextID = 1000
		info.PlayersPerCourt = 4

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

// GetInfo returns the Court class data
func GetInfo() (*Info, error) {

	// Make sure the court directory and info file exists
	_, err := createInfo()
	if err != nil {
		return nil, err
	}

	// Read the court info file
	data, err := ioutil.ReadFile(courtInfoFile)
	if err != nil {
		return nil, err
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, err
	}

	return &i, err
}

// saveInfo save the Court class info
func saveInfo(i Info) error {

	infoJSON, err := json.Marshal(i)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// getAndIncrementCurrentCourtID returns the CurrentID and then increments the CurrentID on disk
func getAndIncrementCurrentCourtID() (int, error) {

	i, err := GetInfo()
	if err != nil {
		return 0, err
	}

	id := i.NextID
	i.NextID = i.NextID + 1

	err = saveInfo(*i)
	if err != nil {
		return 0, err
	}

	return id, nil
}
