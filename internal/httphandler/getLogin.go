package httphandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rsmaxwell/players-api/internal/model"
)

// LogonResponse structure
type LogonResponse struct {
	Token string `json:"token"`
}

// Login method
func Login(rw http.ResponseWriter, req *http.Request) {

	log.Printf("Login:")
	log.Printf("    Method: %s", req.Method)
	log.Printf("    Proto:  %s", req.Proto)
	log.Printf("    Host:   %s", req.Host)
	log.Printf("    URL:    %s", req.URL)

	id, password, _ := req.BasicAuth()

	log.Printf("    id:       %s", id)
	log.Printf("    password: %s", password)

	token, err := model.Login(id, password)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	log.Printf("    token:    %s", token)

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(LogonResponse{
		Token: token,
	})
}
