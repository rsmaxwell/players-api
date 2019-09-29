package common

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"unicode"

	"github.com/rsmaxwell/players-api/codeError"
)

var (
	// RootDir directory
	RootDir string
)

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
		RootDir = home + "/players-api"
	} else {
		RootDir = home + "/players-api-test"
	}
}

// RemoveContents empties the contents of a directory
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return codeError.NewInternalServerError(err.Error())
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return codeError.NewInternalServerError(err.Error())
		}
	}
	return nil
}

// CheckCharactersInID checks the characters are valid for an ID
func CheckCharactersInID(s string) error {
	for _, r := range s {

		ok := false
		if unicode.IsLetter(r) {
			ok = true
		} else if unicode.IsDigit(r) {
			ok = true
		}

		if !ok {
			return codeError.NewBadRequest("Invalid ID")
		}
	}
	return nil
}

// Contains function
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
