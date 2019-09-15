package common

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"unicode"
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

// CheckID checks the characters are valid for an ID
func CheckID(s string) bool {
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
