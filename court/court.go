package court

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
)

// CourtInfo structure
type CourtInfo struct {
	NextID          int `json:"nextID"`
	PlayersPerCourt int `json:"playersPerCourt"`
}

// CreateCourtRequest structure
type CreateCourtRequest struct {
	Token string `json:"token"`
	Court Court  `json:"court"`
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
	logger.Logger.Printf("courtDirectory = %s\n", courtDir)
	logger.Logger.Printf("courtInfoFile = %s\n", courtInfoFile)
}

// CreateCourtDirectory  creates the people directory
func CreateCourtDirectory() error {

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

	return nil
}

// RemoveCourtDirectory removes ALL the court files
func RemoveCourtDirectory() error {
	logger.Logger.Printf("Remove court directory")

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

// CreateCourtInfoFile initialises the Info object from file
func CreateCourtInfoFile() (*CourtInfo, error) {

	logger.Logger.Printf("Checking the courtInfoFile exists")
	if _, err := os.Stat(courtInfoFile); os.IsNotExist(err) {

		err := CreateCourtDirectory()
		if err != nil {
			logger.Logger.Panicf(err.Error())
		}

		info := new(CourtInfo)
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

	var i CourtInfo
	err = json.Unmarshal(data, &i)
	if err != nil {
		logger.Logger.Panicf(err.Error())
	}

	return &i, nil
}

// GetCourtInfo returns the Court class data
func GetCourtInfo() (*CourtInfo, error) {

	// Make sure the court directory and info file exists
	_, err := CreateCourtInfoFile()
	if err != nil {
		return nil, err
	}

	// Read the court info file
	data, err := ioutil.ReadFile(courtInfoFile)
	if err != nil {
		return nil, err
	}

	var i CourtInfo
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, err
	}

	return &i, err
}

// SaveCourtInfo save the Court class info
func SaveCourtInfo(i CourtInfo) error {

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

// GetAndIncrementCurrentCourtID returns the CurrentID and then increments the CurrentID on disk
func GetAndIncrementCurrentCourtID() (int, error) {

	i, err := GetCourtInfo()
	if err != nil {
		return 0, err
	}

	id := i.NextID
	i.NextID = i.NextID + 1

	err = SaveCourtInfo(*i)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// NewCourt initialises a Court object
func NewCourt(name string) (*Court, error) {
	court := new(Court)
	court.Name = name
	return court, nil
}

// UpdateCourt update fields
func UpdateCourt(id int, court2 JSONCourt) (*Court, error) {

	court, err := GetCourtDetails(id)
	if err != nil {
		logger.Logger.Print(err)
		return nil, fmt.Errorf("court [%d] not found", id)
	}

	if court2.Name.Set {
		court.Name = court2.Name.Value
	}

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

	courtfile := courtListDir + "/" + strconv.Itoa(id) + ".json"
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

	files, err := ioutil.ReadDir(courtListDir)
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

	courtfile := courtListDir + "/" + strconv.Itoa(id) + ".json"
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
		logger.Logger.Printf("Could not parse Court data. err = %s\n", err)
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

	courtfile := courtListDir + "/" + strconv.Itoa(id) + ".json"
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
