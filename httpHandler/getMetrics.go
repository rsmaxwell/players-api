package httpHandler

import (
	"encoding/json"
	"net/http"
)

// metrics Response
type metricsResponse struct {
	ClientSuccess             int `json:"clientSuccess"`
	ClientError               int `json:"clientError"`
	ClientAuthenticationError int `json:"clientAuthenticationError"`
	ServerError               int `json:"serverError"`
}

// GetMetrics method
func GetMetrics(rw http.ResponseWriter, req *http.Request) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		WriteResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(metricsResponse{
		ClientSuccess:             clientSuccess,
		ClientError:               clientError,
		ClientAuthenticationError: clientAuthenticationError,
		ServerError:               serverError,
	})
}
