package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rsmaxwell/players-server/logger"
	"github.com/rsmaxwell/players-server/players"
)

var (
	port                      int
	username                  string
	password                  string
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// Error Response JSON object
type messageResponseJSON struct {
	Message string `json:"message"`
}

// metrics Response JSON object
type metricsResponseJSON struct {
	ClientSuccess             int `json:"clientSuccess"`
	ClientError               int `json:"clientError"`
	ClientAuthenticationError int `json:"clientAuthenticationError"`
	ServerError               int `json:"serverError"`
}

// People details Response JSON object
type personDetailsResponseJSON struct {
	Person players.Person `json:"person"`
}

// List people Response JSON object
type listPeopleResponseJSON struct {
	People []int `json:"people"`
}

// Number of people Response JSON object
type numberOfPeopleResponseJSON struct {
	NumberOfPeople int `json:"numberOfPeople"`
}

// Handle writing error messager response.
func writeMessageResponse(w http.ResponseWriter, httpStatus int, message string) {
	logger.Logger.Printf("Response: %d %s - %s", httpStatus, http.StatusText(httpStatus), message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(messageResponseJSON{
		Message: message,
	})
}

// Handle writing numberOfPeople response
func writePersonInfoResponse(w http.ResponseWriter, r *http.Request) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := players.List()
	if err != nil {
		writeMessageResponse(w, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	numberOfPeople := len(listOfPeople)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(numberOfPeopleResponseJSON{
		NumberOfPeople: numberOfPeople,
	})
}

// Handle writing the GET list of People response
func writeGetListOfPeopleResponse(w http.ResponseWriter, r *http.Request) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := players.List()
	if err != nil {
		writeMessageResponse(w, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listPeopleResponseJSON{
		People: listOfPeople,
	})
}

// Write the GET person response
func writeGetPersonDetailsResponse(w http.ResponseWriter, r *http.Request, idString string) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		writeMessageResponse(w, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	person, err := players.Details(id)

	if err != nil {
		writeMessageResponse(w, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(personDetailsResponseJSON{
		Person: *person,
	})
}

// Write the POST Add Person response
func writePostAddPersonResponse(w http.ResponseWriter, r *http.Request) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	var p players.Person

	limitedReader := &io.LimitedReader{R: r.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, fmt.Sprintf("Too much data posted in Add Person request"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &p)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, fmt.Sprintf("Could not parse person data for Add Person request"))
		clientError++
		return
	}

	err = players.AddPerson(p)
	if err != nil {
		writeMessageResponse(w, http.StatusBadRequest, fmt.Sprintf("Could not create a new person. err:%s", err))
		serverError++
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeMessageResponse(w, http.StatusOK, "ok")
}

// Write the DELETE person response
func writeDeletePersonResponse(w http.ResponseWriter, r *http.Request, idString string) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	// Convert the ID into a number
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Logger.Printf(err.Error())
		writeMessageResponse(w, http.StatusBadRequest, fmt.Sprintf("The ID:%s is not a number", idString))
		clientError++
		return
	}

	players.Delete(id)

	w.Header().Set("Content-Type", "application/json")
	writeMessageResponse(w, http.StatusOK, "ok")
}

// Write the GET metrics response
func writeGetMetricsResponse(w http.ResponseWriter, r *http.Request) {
	// Check the user calling the service
	user, pass, _ := r.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(w, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricsResponseJSON{
		ClientSuccess:             clientSuccess,
		ClientError:               clientError,
		ClientAuthenticationError: clientAuthenticationError,
		ServerError:               serverError,
	})
}

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	// PersonInfo
	r.HandleFunc("/personinfo", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writePersonInfoResponse(w, req)
		})).Methods(http.MethodGet)

	// ListPeople
	r.HandleFunc("/people", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetListOfPeopleResponse(w, req)
		})).Methods(http.MethodGet)

	// Person Details
	r.HandleFunc("/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetPersonDetailsResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	// Add Person
	r.HandleFunc("/person", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writePostAddPersonResponse(w, req)
		})).Methods(http.MethodPost)

	// Delete Person
	r.HandleFunc("/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeDeletePersonResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	// Metrics
	r.HandleFunc("/metrics", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetMetricsResponse(w, req)
		})).Methods(http.MethodGet)

	// Not Found
	r.NotFoundHandler = http.HandlerFunc(logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeMessageResponse(w, http.StatusNotFound, "Not Found")
			clientError++
		}))
}

// Simple check on the user calling the service
func checkUser(u, p string) bool {

	fmt.Printf("checkUser: username = %s, %s\n", username, u)
	fmt.Printf("checkUser: password = %s, %s\n", password, p)

	if u == username && p == password {
		return true
	}
	return false
}

func main() {

	logger.Logger.Printf("Players Server")
	var ok bool

	username, ok = os.LookupEnv("USERNAME")
	if !ok {
		username = "foo"
	}

	password, ok = os.LookupEnv("password")
	if !ok {
		password = "foo"
	}

	portstring, ok := os.LookupEnv("port")
	if !ok {
		portstring = "4200"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Registering Router and setting Handlers")
	router := mux.NewRouter()
	setupHandlers(router)

	logger.Logger.Printf("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Success")
}
