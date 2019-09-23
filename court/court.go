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
	"github.com/rsmaxwell/players-api/container"
)

// Info type
type Info struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// Court type
type Court struct {
	Container container.Container `json:"container"`
}

var (
	baseDir  string
	listDir  string
	infoFile string
)

func init() {
	baseDir = common.RootDir + "/courts"
	listDir = baseDir + "/list"
	infoFile = baseDir + "/info.json"
}

// removeAll - remove All the courts
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

// Clear all the courts
func Clear() error {

	err := removeAll()
	if err != nil {
		return err
	}

	_, err = os.Stat(infoFile)
	if err == nil {
		err = os.Remove(infoFile)
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
		return "", err
	}

	filename := listDir + "/" + id + ".json"
	return filename, nil
}

// createDirs creates the people directory
func createDirs() error {

	_, err := os.Stat(listDir)
	if err != nil {
		err := os.MkdirAll(listDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	_, err = GetInfo()
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// New initialises a Court object
func New(name string) *Court {
	court := new(Court)
	court.Container.Name = name
	court.Container.Players = []string{}
	return court
}

// Update method
func Update(id string, court2 map[string]interface{}) (*Court, error) {

	court, err := Load(id)
	if err != nil {
		return nil, codeError.NewNotFound(fmt.Sprintf("Court [%s] not found", id))
	}

	if v, ok := court2["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {
			err = container.Update(&court.Container, container2)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, codeError.NewInternalServerError(fmt.Sprintf("Unexpected Comtainer type: %v", v))
		}
	}

	// Save the updated court to disk
	err = Save(id, court)
	if err != nil {
		return nil, err
	}

	return court, nil
}

// Insert adds a new court to the list
func Insert(court *Court) (string, error) {

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		return "", codeError.NewInternalServerError(err.Error())
	}

	id := strconv.Itoa(count)
	err = Save(id, court)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Save writes a Court to disk
func Save(id string, court *Court) error {

	courtJSON, err := json.Marshal(court)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	filename, err := makeFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// List returns a list of the court IDs
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

// Load returns the Court with the given ID
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

// GetInfo returns the Court class data
func GetInfo() (*Info, error) {

	if _, err := os.Stat(infoFile); os.IsNotExist(err) {

		info := new(Info)
		info.NextID = 1000
		info.PlayersPerCourt = 4

		infoJSON, err := json.Marshal(info)
		if err != nil {
			return nil, codeError.NewInternalServerError(err.Error())
		}

		err = ioutil.WriteFile(infoFile, infoJSON, 0644)
		if err != nil {
			return nil, codeError.NewInternalServerError(err.Error())
		}
	}

	data, err := ioutil.ReadFile(infoFile)
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

// saveInfo save the Court class info
func saveInfo(i Info) error {

	infoJSON, err := json.Marshal(i)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	err = ioutil.WriteFile(infoFile, infoJSON, 0644)
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

	files, err := ioutil.ReadDir(listDir)
	if err != nil {
		return 0, codeError.NewInternalServerError(err.Error())
	}

	return len(files), nil
}

// CheckConsistency function
func CheckConsistency() error {
	return nil
}
