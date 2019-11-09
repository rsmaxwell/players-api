package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/rsmaxwell/players-api/internal/basic/court"
	"github.com/rsmaxwell/players-api/internal/basic/person"
	"github.com/rsmaxwell/players-api/internal/basic/version"
	"github.com/rsmaxwell/players-api/internal/common"

	"github.com/rsmaxwell/players-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	port int
)

func main() {

	log.Printf("Generate-test-data for Players-Api: version: %s\n", version.Version())

	err := createBackupEmpty()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = createBackupOne()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = createBackupLogon()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = createBackupFull()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func createBackupEmpty() error {

	err := createTestdataEmpty()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Backup("empty")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

func createBackupOne() error {

	err := createTestdataOne()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Backup("one")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

func createBackupLogon() error {

	err := createTestdataLoggedin()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Backup("logon")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

func createBackupFull() error {

	err := createTestdataFull()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Startup()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = model.Backup("full")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return nil
}

// Register function
func Register(id, password, firstName, lastName, email string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = person.New(firstName, lastName, email, hashedPassword, false).Save(id)
	if err != nil {
		return err
	}

	return nil
}

// CreateCourt function
func CreateCourt(name string, players []string) error {

	_, err := court.New(name, players).Add()
	if err != nil {
		return err
	}

	return nil
}

func clearModel() error {
	_, err := os.Stat(common.RootDir)
	if err == nil {
		err = removeContents(common.RootDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeContents empties the contents of a directory
func removeContents(dirname string) error {
	children, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, d := range children {
		err = os.RemoveAll(path.Join([]string{dirname, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}
