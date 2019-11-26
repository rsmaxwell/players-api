package court

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
	"github.com/rsmaxwell/players-api/internal/debug"
	"gopkg.in/go-playground/validator.v9"

	"github.com/rsmaxwell/players-api/internal/basic/destination"
	"github.com/rsmaxwell/players-api/internal/basic/peoplecontainer"
)

// Info type
type Info struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// Court type
type Court struct {
	destination.Destination
	Container peoplecontainer.PeopleContainer `json:"container" validate:"required,dive"`
}

var (
	courtBaseDir  string
	courtListDir  string
	courtInfoFile string

	validate *validator.Validate

	pkg = debug.NewPackage("court")

	functionMakeFilename                  = debug.NewFunction(pkg, "makeFilename")
	functionMakefileStructure             = debug.NewFunction(pkg, "makefileStructure")
	functionUpdate                        = debug.NewFunction(pkg, "Update")
	functionUpdateFields                  = debug.NewFunction(pkg, "UpdateFields")
	functionAdd                           = debug.NewFunction(pkg, "Add")
	functionSave                          = debug.NewFunction(pkg, "Save")
	functionList                          = debug.NewFunction(pkg, "List")
	functionExists                        = debug.NewFunction(pkg, "Exists")
	functionLoad                          = debug.NewFunction(pkg, "Load")
	functionRemove                        = debug.NewFunction(pkg, "Remove")
	functionGetCourtInfo                  = debug.NewFunction(pkg, "GetCourtInfo")
	functionSaveCourtInfo                 = debug.NewFunction(pkg, "saveCourtInfo")
	functionGetAndIncrementCurrentCourtID = debug.NewFunction(pkg, "getAndIncrementCurrentCourtID")
	functionSize                          = debug.NewFunction(pkg, "Size")
	functionCheckSpace                    = debug.NewFunction(pkg, "CheckSpace")
)

func init() {
	courtBaseDir = common.RootDir + "/courts"
	courtListDir = courtBaseDir + "/list"
	courtInfoFile = courtBaseDir + "/info.json"

	validate = validator.New()
}

// makeFilename function
func makeFilename(id string) (string, error) {
	f := functionMakeFilename

	err := common.CheckCharactersInID(id)
	if err != nil {
		message := fmt.Sprintf("charactor check failed for court id [%s]: %v", id, err)
		f.Dump(message)
		return "", codeerror.NewInternalServerError(message)
	}

	err = makefileStructure()
	if err != nil {
		message := fmt.Sprintf("could not make court file structure: %v", err)
		f.Dump(message)
		return "", codeerror.NewInternalServerError(message)
	}

	filename := courtListDir + "/" + id + ".json"
	return filename, nil
}

// makefileStructure creates the people directory
func makefileStructure() error {
	f := functionMakefileStructure

	_, err := os.Stat(courtListDir)
	if err != nil {
		err := os.MkdirAll(courtListDir, 0755)
		if err != nil {
			message := fmt.Sprintf("could not make the directory [%s]: %v", courtListDir, err)
			f.Dump(message)
			return codeerror.NewInternalServerError(message)
		}
	}

	return nil
}

// New initialises a Court object
func New(name string, players []string) *Court {
	c := new(Court)
	c.Container.Name = name
	c.Container.Players = players
	return c
}

// Update method
func Update(ref *common.Reference, fields map[string]interface{}) error {
	f := functionUpdate
	f.DebugVerbose("ref: %v, fields: %v:", ref, fields)

	c, err := Load(ref)
	if err != nil {
		return err
	}

	err = c.updateFields(fields)
	if err != nil {
		return err
	}

	err = c.Save(ref)
	if err != nil {
		message := fmt.Sprintf("could not save the court [%v]: %v", ref, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}

// updateFields method
func (c *Court) updateFields(fields map[string]interface{}) error {
	f := functionUpdateFields

	if v, ok := fields["Container"]; ok {
		if container2, ok := v.(map[string]interface{}); ok {

			err := c.Container.Update(container2)
			if err != nil {
				return err
			}
		} else {
			message := fmt.Sprintf("Unexpected Container type: %t   %v", v, v)
			f.Dump(message)
			return codeerror.NewInternalServerError(message)
		}
	}

	return nil
}

// Add a new court to the list
func (c *Court) Add() (string, error) {
	f := functionAdd

	count, err := getAndIncrementCurrentCourtID()
	if err != nil {
		message := fmt.Sprintf("could not increment the court counter: %v", err)
		f.Dump(message)
		return "", codeerror.NewInternalServerError(message)
	}

	id := strconv.Itoa(count)
	ref := common.Reference{Type: "court", ID: id}
	err = c.Save(&ref)
	if err != nil {
		message := fmt.Sprintf("could save the court[%s]: %v", id, err)
		f.Dump(message)
		return "", codeerror.NewInternalServerError(message)
	}

	f.DebugVerbose("id: %s", id)
	return id, nil
}

// Save writes a Court to disk
func (c *Court) Save(ref *common.Reference) error {
	f := functionSave

	if ref.Type != "court" {
		message := fmt.Sprintf("Unexpected Reference type[%s]", ref.Type)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	err := validate.Struct(c)
	if err != nil {
		message := fmt.Sprintf("validation for court failed: %v", err)
		f.Dump(message)
		return codeerror.NewBadRequest(message)
	}

	courtJSON, err := json.Marshal(c)
	if err != nil {
		message := fmt.Sprintf("could not marshal court: %v", err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	filename, err := makeFilename(ref.ID)
	if err != nil {
		message := fmt.Sprintf("could not make filename for court[%s]: %v", ref.ID, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	err = ioutil.WriteFile(filename, courtJSON, 0644)
	if err != nil {
		message := fmt.Sprintf("could not write file [%s] for court[%s]: %v", filename, ref.ID, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}

// List returns a list of the court IDs
func List() ([]string, error) {
	f := functionList

	err := makefileStructure()
	if err != nil {
		f.Dump("could not make court file structure: %v", err)
		return nil, err
	}

	files, err := ioutil.ReadDir(courtListDir)
	if err != nil {
		message := fmt.Sprintf("could not read the directory[%s]: %v", courtListDir, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
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
	f := functionExists

	filename, err := makeFilename(id)
	if err != nil {
		f.Dump("could not make filename for court[%s]: %v", id, err)
		return false
	}

	_, err = os.Stat(filename)
	if err != nil {
		f.DebugVerbose("could not stat filename[%s] for court[%s]: %v", filename, id, err)
		return false
	}

	return true
}

// Load returns the Court with the given ID
func Load(ref *common.Reference) (*Court, error) {
	f := functionLoad

	if ref.Type != "court" {
		message := fmt.Sprintf("Unexpected Reference type for court[%v]", ref)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	filename, err := makeFilename(ref.ID)
	if err != nil {
		f.Dump("could not make filename for court[%s]", ref.ID)
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeerror.NewNotFound(fmt.Sprintf("File not found: [%s]", filename))
		}
		message := fmt.Sprintf("could not read court[%s] file [%s]: %v", ref.Type, filename, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	var c Court
	err = json.Unmarshal(data, &c)
	if err != nil {
		message := fmt.Sprintf("could not unmarshal court [%s]", filename)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}
	return &c, nil
}

// Remove the court with the given ID
func Remove(id string) error {
	f := functionRemove

	filename, err := makeFilename(id)
	if err != nil {
		f.Dump("could not make file for court [%s]", id)
		return err
	}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return codeerror.NewNotFound(fmt.Sprintf("file not found: [%s]", filename))
		} else {
			message := fmt.Sprintf("could not stat: [%s]: %v", filename, err)
			f.Dump(message)
			return codeerror.NewInternalServerError(message)
		}
	}

	err = os.Remove(filename)
	if err != nil {
		message := fmt.Sprintf("could not remove the file [%s] for court [%s]", filename, id)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}
	return nil
}

// GetCourtInfo returns the Court class data
func GetCourtInfo() (*Info, error) {
	f := functionGetCourtInfo

	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		info := new(Info)
		info.NextID = 1000
		info.PlayersPerCourt = 4

		infoJSON, err := json.Marshal(info)
		if err != nil {
			message := fmt.Sprintf("could not marshal court info: %v", err)
			f.Dump(message)
			return nil, codeerror.NewInternalServerError(message)
		}

		err = makefileStructure()
		if err != nil {
			f.Dump("could not make the court file structure: %v", err)
			return nil, err
		}

		err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
		if err != nil {
			message := fmt.Sprintf("could not write the court info file [%s]: %v", courtInfoFile, err)
			f.Dump(message)
			return nil, codeerror.NewInternalServerError(message)
		}
	}

	data, err := ioutil.ReadFile(courtInfoFile)
	if err != nil {
		message := fmt.Sprintf("could not read read the court info file [%s]: %v", courtInfoFile, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	var i Info
	err = json.Unmarshal(data, &i)
	if err != nil {
		message := fmt.Sprintf("could not unmarshal the court info file [%s]: %v", courtInfoFile, err)
		f.Dump(message)
		return nil, codeerror.NewInternalServerError(message)
	}

	return &i, nil
}

// saveCourtInfo save the Court class info
func saveCourtInfo(i Info) error {
	f := functionSaveCourtInfo

	infoJSON, err := json.Marshal(i)
	if err != nil {
		message := fmt.Sprintf("could not marshal the court info: %v", err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	err = makefileStructure()
	if err != nil {
		f.Dump("could not make the court file structure: %v", err)
		return err
	}

	err = ioutil.WriteFile(courtInfoFile, infoJSON, 0644)
	if err != nil {
		message := fmt.Sprintf("could not write the court file[%s]: %v", courtInfoFile, err)
		f.Dump(message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}

// getAndIncrementCurrentCourtID returns the CurrentID and then increments the CurrentID on disk
func getAndIncrementCurrentCourtID() (int, error) {
	f := functionGetAndIncrementCurrentCourtID

	i, err := GetCourtInfo()
	if err != nil {
		f.Dump("could not get the court info: %v", err)
		return 0, err
	}

	id := i.NextID
	i.NextID = i.NextID + 1

	err = saveCourtInfo(*i)
	if err != nil {
		f.Dump("could not save the court info: %v", err)
		return 0, err
	}

	return id, nil
}

// Size returns the number of courts
func Size() (int, error) {
	f := functionSize

	files, err := ioutil.ReadDir(courtListDir)
	if err != nil {
		message := fmt.Sprintf("could not read the court list directory[%s]: %v", courtListDir, err)
		f.Dump(message)
		return 0, codeerror.NewInternalServerError(message)
	}

	return len(files), nil
}

// GetContainer returns the Destination
func (c *Court) GetContainer() *peoplecontainer.PeopleContainer {
	return &c.Container
}

// CheckPlayersLocation checks the players are at this destination
func (c *Court) CheckPlayersLocation(players []string) error {
	pc := c.GetContainer()
	return destination.CheckPlayersInContainer(pc, players)
}

// CheckSpace checks there is space in the target for the moving players
func (c *Court) CheckSpace(players []string) error {
	f := functionCheckSpace

	info, err := GetCourtInfo()
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	containerSize := len(c.Container.Players)
	playersSize := len(players)

	if containerSize+playersSize > info.PlayersPerCourt {
		message := fmt.Sprintf("Too many players. %d + %d > %d", containerSize, playersSize, info.PlayersPerCourt)
		f.Dump(message)
		return codeerror.NewBadRequest(message)
	}

	return nil
}

// RemovePlayers deletes players from the destination
func (c *Court) RemovePlayers(players []string) error {
	pc := c.GetContainer()
	return destination.RemovePlayersFromContainer(pc, players)
}

// AddPlayers adds players to the destination
func (c *Court) AddPlayers(players []string) error {
	pc := c.GetContainer()
	return destination.AddPlayersToContainer(pc, players)
}
