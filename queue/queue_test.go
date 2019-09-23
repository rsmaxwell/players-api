package queue

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/stretchr/testify/require"
)

// RemoveDir - Remove the Queue file
func RemoveDir() error {

	filename, err := makeFilename()
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err == nil {
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestClearQueue(t *testing.T) {
	r := require.New(t)

	queuefile, err := makeFilename()
	r.Nil(err, "err should be nothing")

	err = Clear()
	r.Nil(err, "err should be nothing")

	_, err = os.Stat(queuefile)
	r.Nil(err, "err should be nothing")

	err = RemoveDir()
	r.Nil(err, "err should be nothing")

	_, err = os.Stat(queuefile)
	if err != nil {
		if os.IsNotExist(err) { // File does not exist
			// ok
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: %s", err.Error()))
		}
	} else {
		r.Fail("The queue file was not removed")
	}

}

func TestNewInfoJunkQueue(t *testing.T) {
	r := require.New(t)

	err := Clear()
	r.Nil(err, "err should be nothing")

	queuefile, err := makeFilename()
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(queuefile, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	_, err = Load()
	if err != nil {
		if cerr, ok := err.(*codeError.CodeError); ok {
			if cerr.Code() == http.StatusInternalServerError {
				// ok
			} else {
				r.Fail(fmt.Sprintf("Unexpected error: %d", cerr.Code()))
			}
		} else {
			r.Fail(fmt.Sprintf("Unexpected error: %s", err.Error()))
		}
	} else {
		r.Fail("Unexpected success")
	}
}
