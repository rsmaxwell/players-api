package model

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/sync"

	"github.com/stretchr/testify/require"
)

// BasicAuth function
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// SetupEmpty function
func SetupEmpty(t *testing.T) func(t *testing.T) {

	err := Restore("empty")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupOne function
func SetupOne(t *testing.T) func(t *testing.T) {

	err := Restore("one")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupLoggedin function
func SetupLoggedin(t *testing.T) func(t *testing.T) {

	err := Restore("logon")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupFull function
func SetupFull(t *testing.T) func(t *testing.T) {

	err := Restore("full")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

func listdir(title string, root string) error {

	log.Printf("%s: %s\n", title, root)

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range fileInfo {
		log.Printf("    %t  %o  %s\n", file.IsDir(), file.Mode(), file.Name())
	}
	return nil
}

// Backup function
func Backup(name string) error {

	reference := common.RootDir
	copy := common.RootDir + "-backup/" + name

	listdir("Backup: before:", filepath.Dir(common.RootDir))

	err := os.MkdirAll(copy, 0755)
	if err != nil {
		return err
	}

	log.Printf("Backup: after\n")

	err = sync.Dir(reference, copy)
	if err != nil {
		return err
	}

	listdir("Backup: after:", filepath.Dir(common.RootDir))

	return nil
}

// Restore function
func Restore(name string) error {

	reference := common.RootDir + "-backup/" + name
	copy := common.RootDir

	err := sync.Dir(reference, copy)
	if err != nil {
		return err
	}

	return nil
}
