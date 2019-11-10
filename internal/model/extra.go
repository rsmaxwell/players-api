package model

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/sync"

	"github.com/stretchr/testify/require"
)

var (
	functionBasicAuth     = debug.NewFunction(pkg, "BasicAuth")
	functionSetupEmpty    = debug.NewFunction(pkg, "SetupEmpty")
	functionSetupOne      = debug.NewFunction(pkg, "SetupOne")
	functionSetupLoggedin = debug.NewFunction(pkg, "SetupLoggedin")
	functionSetupFull     = debug.NewFunction(pkg, "SetupFull")
	functionListdir       = debug.NewFunction(pkg, "listdir")
	functionBackup        = debug.NewFunction(pkg, "Backup")
	functionRestore       = debug.NewFunction(pkg, "Restore")
)

// BasicAuth function
func BasicAuth(username, password string) string {
	f := functionBasicAuth
	f.DebugInfo("username: [%s], password:[%s]", username, password)

	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// SetupEmpty function
func SetupEmpty(t *testing.T) func(t *testing.T) {
	f := functionSetupEmpty
	f.DebugInfo("")

	err := Restore("empty")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupOne function
func SetupOne(t *testing.T) func(t *testing.T) {
	f := functionSetupOne
	f.DebugInfo("")

	err := Restore("one")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupLoggedin function
func SetupLoggedin(t *testing.T) func(t *testing.T) {
	f := functionSetupLoggedin
	f.DebugInfo("")

	err := Restore("logon")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

// SetupFull function
func SetupFull(t *testing.T) func(t *testing.T) {
	f := functionSetupFull
	f.DebugInfo("")

	err := Restore("full")
	require.Nil(t, err, "err should be nothing")

	return func(t *testing.T) {
	}
}

func listdir(title string, root string) error {
	f := functionListdir
	f.DebugInfo("%s: %s", title, root)

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		f.Dump("could not make the read the root directory [%s]\n%v", root, err)
		return err
	}
	for _, file := range fileInfo {
		f.DebugInfo("    %t  %o  %s\n", file.IsDir(), file.Mode(), file.Name())
	}
	return nil
}

// Backup function
func Backup(name string) error {
	f := functionBackup
	f.DebugInfo("name: %s", name)

	reference := common.RootDir
	copy := common.RootDir + "-backup/" + name

	listdir("Backup(1)", filepath.Dir(common.RootDir))

	err := os.MkdirAll(copy, 0755)
	if err != nil {
		f.Dump("could not make the copy directory [%s]\n%v", copy, err)
		return err
	}

	listdir("Backup(2)", filepath.Dir(common.RootDir))

	err = sync.HandleDir(reference, copy)
	if err != nil {
		f.Dump("could not sync [%s] with [%s]\n%v", reference, copy, err)
		return err
	}

	listdir("Backup(3)", filepath.Dir(common.RootDir))

	return nil
}

// Restore function
func Restore(name string) error {
	f := functionRestore
	f.DebugInfo("name: %s", name)

	reference := common.RootDir + "-backup/" + name
	copy := common.RootDir

	err := sync.HandleDir(reference, copy)
	if err != nil {
		f.Dump("could not sync [%s] with [%s]\n%v", reference, copy, err)
		return err
	}

	return nil
}
