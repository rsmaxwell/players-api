package queue

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rsmaxwell/players-api/codeError"
	"github.com/stretchr/testify/require"
)

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
