package destination

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
)

// Info type
type Info struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// Court type
type Court struct {
	Destination
	Container Container `json:"container"`
}

var (
	listDir  string
	infoFile string
)

func init() {
	baseDir = common.RootDir + "/courts"
	listDir = baseDir + "/list"
	infoFile = baseDir + "/info.json"
}

// removeAll - remove All the courts
func removeAllCourts() error {

	_, err := os.Stat(listDir)
	if err == nil {
		err = common.RemoveContents(listDir)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	return nil
}

// ClearCourts all the courts
func ClearCourts() error {

	err := removeAllCourts()
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

	err = createCourtFiles()
	if err != nil {
		return err
	}

	return nil
}

// makeCourtFilename function
func makeCourtFilename(id string) (string, error) {

	err := common.CheckCharactersInID(id)
	if err != nil {
		return "", err
	}

	filename := listDir + "/" + id + ".json"
	return filename, nil
}

// createDirs creates the people directory
func createCourtFiles() error {

	_, err := os.Stat(listDir)
	if err != nil {
		err := os.MkdirAll(listDir, 0755)
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}

	_, err = GetCourtInfo()
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// NewCourt initialises a Court object
func NewCourt(name string) *Court {
	c := new(Court)
	c.Container.Name = name
	c.Container.Players = []string{}
	return c
}

// UpdateCourt method
func UpdateCourt(id string, fields map[string]interface{}) error {

	c, err := LoadCourt(id)
	if err != nil {
		return err
	}

	if v, ok := fields["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {
			err := c.Container.Update(container2)
			if err != nil {
				return err
			}
		} else {
			return codeError.NewInternalServerError(fmt.Sprintf("Unexpected Container type: %v", v))
		}
	}

	err = c.Save(id)
	if err != nil {
		return err
	}

	return nil
}

// Insert adds a new court to the list
func (c Court) Insert() (string, error) {

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		return "", codeError.NewInternalServerError(err.Error())
	}

	id := strconv.Itoa(count)
	err = c.Save(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Save writes a Court to disk
func (c Court) Save(id string) error {

	courtJSON, err := json.Marshal(c)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	filename, err := makeCourtFilename(id)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}

	return nil
}

// ListCourts returns a list of the court IDs
func ListCourts() ([]string, error) {

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

// CourtExists returns 'true' if the court exists
func CourtExists(id string) bool {

	filename, err := makeCourtFilename(id)
	if err != nil {
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

// LoadCourt returns the Court with the given ID
func LoadCourt(id string) (*Court, error) {

	filename, err := makeCourtFilename(id)
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

// RemoveCourt the court with the given ID
func RemoveCourt(id string) error {

	filename, err := makeCourtFilename(id)
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

// GetCourtInfo returns the Court class data
func GetCourtInfo() (*Info, error) {

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

// saveCourtInfo save the Court class info
func saveCourtInfo(i Info) error {

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

	i, err := GetCourtInfo()
	if err != nil {
		return 0, err
	}

	id := i.NextID
	i.NextID = i.NextID + 1

	err = saveCourtInfo(*i)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CourtSize returns the number of courts
func CourtSize() (int, error) {

	files, err := ioutil.ReadDir(listDir)
	if err != nil {
		return 0, codeError.NewInternalServerError(err.Error())
	}

	return len(files), nil
}

// GetContainer returns the Destination
func (c Court) GetContainer() *Container {
	return &c.Container
}
