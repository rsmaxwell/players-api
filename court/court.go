package court

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/common"
	"github.com/rsmaxwell/players-api/person"
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

// removeAllCourts - remove the list of courts
func removeAllCourts() error {

	_, err := os.Stat(courtListDir)
	if err == nil {
		err = common.RemoveContents(courtListDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// Clear - Clear the list of courts
func Clear() error {

	err := removeAllCourts()
	if err != nil {
		return err
	}

	_, err = os.Stat(courtInfoFile)
	if err == nil {
		err = os.Remove(courtInfoFile)
		if err != nil {
			return err
		}
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
		return "", codeError.NewBadRequest(err.Error())
	}

	filename := courtListDir + "/" + id + ".json"
	return filename, nil
}

// createDirs creates the people directory
func createDirs() error {

	_, err := os.Stat(courtListDir)
	if err != nil {
		err := os.MkdirAll(courtListDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	_, err = createInfo()
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
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
func Update(id string, court2 map[string]interface{}) (*Court, error) {

	court, err := Load(id)
	if err != nil {
		return nil, codeError.NewNotFound(fmt.Sprintf("Court [%s] not found", id))
	}

	if v, ok := court2["Name"]; ok {
		value, ok := v.(string)
		if !ok {
			return nil, codeError.NewBadRequest(fmt.Sprintf("'Name' was an unexpected type. Expected: 'string'. Actual: %T", v))
		}
		court.Name = value
	}

	if v, ok := court2["Players"]; ok {

		array := []string{}
		for _, i := range v.([]interface{}) {

			id2, ok := i.(string)
			if !ok {
				return nil, codeError.NewBadRequest(fmt.Sprintf("'Player' was an unexpected type. Expected: 'string'. Actual: %T", i))
			}

			if !person.Exists(id2) {
				return nil, codeError.NewNotFound(fmt.Sprintf("Person [%s] not found", id2))
			}

			if !person.IsPlayer(id2) {
				return nil, codeError.NewBadRequest(fmt.Sprintf("Person [%s] is not a player", id2))
			}

			array = append(array, id2)
		}
		court.Players = array
	}

	// Convert the 'court' object into a JSON string
	courtJSON, err := json.Marshal(court)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	// Save the updated court to disk
	filename, err := makeFilename(id)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	return court, nil
}

// Save adds a court to the list of courts
func Save(court *Court) (string, error) {

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		return "", codeError.NewInternalServerError(err.Error())
	}

	courtJSON, err := json.Marshal(court)
	if err != nil {
		return "", codeError.NewInternalServerError(err.Error())
	}

	id := strconv.Itoa(count)
	filename := courtListDir + "/" + id + ".json"
	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		return "", codeError.NewInternalServerError(err.Error())
	}

	return id, nil
}

// List returns a list of the court IDs
func List() ([]string, error) {

	files, err := ioutil.ReadDir(courtListDir)
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

// Exists returns 'true' if the court exists
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

// Get returns the details of the court with the given ID
func Load(id string) (*Court, error) {

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

	var c Court
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}
	return &c, nil
}

// Remove the court with the given ID
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

// createInfo initialises the Info objects
func createInfo() (*Info, error) {

	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		info := new(Info)
		info.NextID = 1000
		info.PlayersPerCourt = 4

		infoJSON, err := json.Marshal(info)
		if err != nil {
			return nil, codeError.NewInternalServerError(err.Error())
		}

		err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
		if err != nil {
			return nil, codeError.NewInternalServerError(err.Error())
		}
	}

	data, err := ioutil.ReadFile(courtInfoFile)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
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
		return nil, codeError.NewInternalServerError(err.Error())
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, codeError.NewInternalServerError(err.Error())
	}

	return &i, err
}

// saveInfo save the Court class info
func saveInfo(i Info) error {

	infoJSON, err := json.Marshal(i)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
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

// Size returns the number of courts
func Size() (int, error) {

	files, err := ioutil.ReadDir(courtListDir)
	if err != nil {
		return 0, codeError.NewInternalServerError(err.Error())
	}

	return len(files), nil
}
