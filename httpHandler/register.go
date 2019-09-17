package httpHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rsmaxwell/players-api/person"
	"golang.org/x/crypto/bcrypt"
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

	var r RegisterRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		WriteResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse Person data"))
		clientError++
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, fmt.Sprintf("Could not generate hash"))
		serverError++
	}

	p, err := person.New(r.FirstName, r.LastName, r.Email, hashedPassword, false)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, fmt.Sprintf("Could not create new Person"))
		serverError++
	}

	err = person.Add(r.UserID, *p)
	if err != nil {
		WriteResponse(rw, http.StatusInternalServerError, fmt.Sprintf("Could not add Person"))
		serverError++
	}

	setHeaders(rw, req)
	WriteResponse(rw, http.StatusOK, "ok")
}
