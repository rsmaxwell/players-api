package players

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"unicode"

	"github.com/rsmaxwell/players-api/logger"
)

var (
	rootdir string
)

// Info structure
type Info struct {
	CurrentID int `json:"currentID"`
}

func homeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func init() {

	home := homeDir()

	if flag.Lookup("test.v") == nil {
		rootdir = home + "/players-api"
	} else {
		rootdir = home + "/players-api-test"
	}

	peopleDirectory = rootdir + "/people"
	peopleDataDirectory = peopleDirectory + "/data"
	logger.Logger.Printf("peopleDirectory = %s\n", peopleDirectory)

	courtDirectory = rootdir + "/court"
	courtDataDirectory = courtDirectory + "/data"
	logger.Logger.Printf("courtDirectory = %s\n", courtDirectory)
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func checkID(s string) bool {
	for _, r := range s {

		ok := false
		if unicode.IsLetter(r) {
			ok = true
		} else if unicode.IsDigit(r) {
			ok = true
		}

		if !ok {
			return false
		}
	}
	return true
}