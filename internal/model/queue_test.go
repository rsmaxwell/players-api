package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"

	"github.com/stretchr/testify/require"
)

func TestNewInfoJunkQueue(t *testing.T) {
	r := require.New(t)

	err := ClearQueue()
	r.Nil(err, "err should be nothing")

	queuefile, err := makeQueueFilename()
	r.Nil(err, "err should be nothing")

	err = ioutil.WriteFile(queuefile, []byte("junk"), 0644)
	r.Nil(err, "err should be nothing")

	ref := common.Reference{Type: "queue", ID: ""}
	_, err = LoadQueue(&ref)
	if err != nil {
		if cerr, ok := err.(*codeerror.CodeError); ok {
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
