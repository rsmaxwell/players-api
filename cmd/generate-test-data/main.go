package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/rsmaxwell/players-api/internal/common"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/model"
	"github.com/rsmaxwell/players-api/internal/session"
	"golang.org/x/crypto/bcrypt"
)

var (
	port int
)

func main() {

	log.Printf("Generate-test-data for Players Server: 2019-10-03 08:55\n")

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

	err := createTestdataLoggedon()
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

	err = model.NewPerson(firstName, lastName, email, hashedPassword, false).Save(id)
	if err != nil {
		return err
	}

	return nil
}

// Login function
func Login(user, pass string) (string, error) {

	if !model.CheckUser(user, pass) {
		return "", errors.New("Invalid Userid or Password")
	}

	token, err := session.New(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateCourt function
func CreateCourt(name string, players []string) error {

	_, err := model.NewCourt(name, players).Add()
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
func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return codeerror.NewInternalServerError(err.Error())
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return codeerror.NewInternalServerError(err.Error())
		}
	}
	return nil
}
