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
	"github.com/rsmaxwell/players-api/internal/debug"

	"github.com/rsmaxwell/players-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	pkg = debug.NewPackage("main")

	functionMain              = debug.NewFunction(pkg, "main")
	functionCreateBackupEmpty = debug.NewFunction(pkg, "createBackupEmpty")
	functionCreateBackupOne   = debug.NewFunction(pkg, "createBackupOne")
	functionCreateBackupLogon = debug.NewFunction(pkg, "createBackupLogon")
	functionCreateBackupFull  = debug.NewFunction(pkg, "createBackupFull")
	functionRegister          = debug.NewFunction(pkg, "Register")
	functionCreateCourt       = debug.NewFunction(pkg, "CreateCourt")
	functionClearModel        = debug.NewFunction(pkg, "clearModel")
	functionRemoveContents    = debug.NewFunction(pkg, "removeContents")
)

var (
	port int
)

func main() {
	f := functionMain
	f.Infof("Generate-test-data for Players-Api: version: %s\n", version.Version())

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
	f := functionCreateBackupEmpty
	f.DebugInfo("\n")

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
	f := functionCreateBackupOne
	f.DebugInfo("\n")

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
	f := functionCreateBackupLogon
	f.DebugInfo("\n")

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
	f := functionCreateBackupFull
	f.DebugInfo("\n")

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
	f := functionRegister
	f.DebugInfo("id: %s, password: %s, firstName: %s, lastName: %s, email: %s\n", id, password, firstName, lastName, email)

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
	f := functionCreateCourt
	f.DebugInfo("name: %s, players: %v\n", name, players)

	_, err := court.New(name, players).Add()
	if err != nil {
		return err
	}

	return nil
}

func clearModel() error {
	f := functionClearModel
	f.DebugInfo("\n")

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
	f := functionRemoveContents
	f.DebugInfo("dirname: %s\n", dirname)

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
