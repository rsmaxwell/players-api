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
	"golang.org/x/crypto/bcrypt"

	"github.com/rsmaxwell/players-api/internal/model"
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
	f.DebugVerbose("")

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
	f.DebugVerbose("")

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
	f.DebugVerbose("")

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
	f.DebugVerbose("")

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
	f.DebugVerbose("id: %s, password: %s, firstName: %s, lastName: %s, email: %s", id, password, firstName, lastName, email)

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
	f.DebugVerbose("name: %s, players: %v", name, players)

	_, err := court.New(name, players).Add()
	if err != nil {
		return err
	}

	return nil
}

func clearModel() error {
	f := functionClearModel
	f.DebugVerbose("")

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
	f.DebugVerbose("dirname: %s", dirname)

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
