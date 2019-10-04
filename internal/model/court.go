package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/session"
)

// Info type
type Info struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// Court type
type Court struct {
	Destination
	Container PeopleContainer `json:"container" validate:"required,dive"`
}

var (
	courtBaseDir  string
	courtListDir  string
	courtInfoFile string
)

func init() {
	courtBaseDir = common.RootDir + "/courts"
	courtListDir = courtBaseDir + "/list"
	courtInfoFile = courtBaseDir + "/info.json"
}

// removeAll - remove All the courts
func removeAllCourts() error {

	_, err := os.Stat(courtListDir)
	if err == nil {
		err = common.RemoveContents(courtListDir)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
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

	_, err = os.Stat(courtInfoFile)
	if err == nil {
		err = os.Remove(courtInfoFile)
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

	filename := courtListDir + "/" + id + ".json"
	return filename, nil
}

// createDirs creates the people directory
func createCourtFiles() error {

	_, err := os.Stat(courtListDir)
	if err != nil {
		err := os.MkdirAll(courtListDir, 0755)
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}

	_, err = GetCourtInfo()
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	return nil
}

// NewCourt initialises a Court object
func NewCourt(name string, players []string) *Court {
	c := new(Court)
	c.Container.Name = name
	c.Container.Players = players
	return c
}

// UpdateCourt method
func UpdateCourt(ref *common.Reference, session *session.Session, fields map[string]interface{}) error {

	c, err := LoadCourt(ref)
	if err != nil {
		return err
	}

	err = c.Update(session, fields)
	if err != nil {
		return err
	}

	err = c.Save(ref)
	if err != nil {
		return err
	}

	return nil
}

// Update method
func (c *Court) Update(session *session.Session, fields map[string]interface{}) error {

	if v, ok := fields["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {

			err := c.Container.Update(container2)
			if err != nil {
				return err
			}
		} else {
			return codeerror.NewInternalServerError(fmt.Sprintf("Unexpected Container type: %t   %v", v, v))
		}
	}

	return nil
}

// Add a new court to the list
func (c *Court) Add() (string, error) {

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		return "", codeerror.NewInternalServerError(err.Error())
	}

	id := strconv.Itoa(count)
	ref := common.Reference{Type: "court", ID: id}
	err = c.Save(&ref)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Save writes a Court to disk
func (c *Court) Save(ref *common.Reference) error {

	if ref.Type != "court" {
		return codeerror.NewInternalServerError("Unexpected Reference type")
	}

	err := validate.Struct(c)
	if err != nil {
		return codeerror.NewBadRequest(err.Error())
	}

	courtJSON, err := json.Marshal(c)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	filename, err := makeCourtFilename(ref.ID)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	return nil
}

// ListCourts returns a list of the court IDs
func ListCourts() ([]string, error) {

	files, err := ioutil.ReadDir(courtListDir)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
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
func LoadCourt(ref *common.Reference) (*Court, error) {

	if ref.Type != "court" {
		return nil, codeerror.NewInternalServerError("Unexpected Reference type")
	}

	filename, err := makeCourtFilename(ref.ID)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeerror.NewNotFound(err.Error())
		}
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	var c Court
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
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
			return codeerror.NewInternalServerError(err.Error())
		}
		return nil

	} else if os.IsNotExist(err) { // File does not exist
		return codeerror.NewNotFound(fmt.Sprintf("File Not Found: %s", filename))
	}

	return codeerror.NewInternalServerError(err.Error())
}

// GetCourtInfo returns the Court class data
func GetCourtInfo() (*Info, error) {

	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		info := new(Info)
		info.NextID = 1000
		info.PlayersPerCourt = 4

		infoJSON, err := json.Marshal(info)
		if err != nil {
			return nil, codeerror.NewInternalServerError(err.Error())
		}

		err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
		if err != nil {
			return nil, codeerror.NewInternalServerError(err.Error())
		}
	}

	data, err := ioutil.ReadFile(courtInfoFile)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	return &i, nil
}

// saveCourtInfo save the Court class info
func saveCourtInfo(i Info) error {

	infoJSON, err := json.Marshal(i)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}

	err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
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

	files, err := ioutil.ReadDir(courtListDir)
	if err != nil {
		return 0, codeerror.NewInternalServerError(err.Error())
	}

	return len(files), nil
}

// GetContainer returns the Destination
func (c *Court) GetContainer() *PeopleContainer {
	return &c.Container
}

// CheckPlayersLocation checks the players are at this destination
func (c *Court) CheckPlayersLocation(players []string) error {
	pc := c.GetContainer()
	return CheckPlayersInContainer(pc, players)
}

// CheckSpace checks there is space in the target for the moving players
func (c *Court) CheckSpace(players []string) error {
	info, err := GetCourtInfo()
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	containerSize := len(c.Container.Players)
	playersSize := len(players)

	if containerSize+playersSize > info.PlayersPerCourt {
		return codeerror.NewBadRequest(fmt.Sprintf("Too many players. %d + %d > %d", containerSize, playersSize, info.PlayersPerCourt))
	}

	return nil
}

// RemovePlayers deletes players from the destination
func (c *Court) RemovePlayers(players []string) error {
	pc := c.GetContainer()
	return RemovePlayersFromContainer(pc, players)
}

// AddPlayers adds players to the destination
func (c *Court) AddPlayers(players []string) error {
	pc := c.GetContainer()
	return AddPlayersToContainer(pc, players)
}
