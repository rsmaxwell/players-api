package httphandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/model"
)

// RegisterRequest structure
type RegisterRequest struct {
	UserID    string `json:"userID"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

// Register method
func Register(rw http.ResponseWriter, req *http.Request) {

	log.Printf("Register:")
	log.Printf("    Method: %s", req.Method)
	log.Printf("    Proto:  %s", req.Proto)
	log.Printf("    Host:   %s", req.Host)
	log.Printf("    URL:    %s", req.URL)

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	var r RegisterRequest
	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, err.Error())
		common.MetricsData.ClientError++
		return
	}

	log.Printf("    UserID:    %s", r.UserID)
	log.Printf("    Password:  %s", r.Password)
	log.Printf("    FirstName: %s", r.FirstName)
	log.Printf("    LastName:  %s", r.LastName)
	log.Printf("    Email:     %s", r.Email)

	err = model.Register(r.UserID, r.Password, r.FirstName, r.LastName, r.Email)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
