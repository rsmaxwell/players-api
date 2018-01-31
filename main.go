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

	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/players"
)

var (
	port                      int
	username                  string
	password                  string
	baseURL                   string
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// Authenticate Response JSON object
type authenticateResponseJSON struct {
	Token string `json:"token"`
}

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

func setHeaders(rw http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")

	if origin == "" {
		origin = "http://localhost:4200"
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", origin)
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Authorization")
}

// Handle authenticate response
func writeAuthenticateResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeAuthenticateResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	token := "qwerty"

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(authenticateResponseJSON{
		Token: token,
	})
}

// Handle Users - getAll
func writeUsersGetAllResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeUsersgetAllResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Handle Users - get by Id
func writeUsersGetByIDResponse(rw http.ResponseWriter, req *http.Request, id string) {

	logger.Logger.Printf("writeUsersGetByIdResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Handle Users - create
func writeUsersCreateResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeUsersCreateResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Handle Users - update
func writeUsersUpdateResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeUsersUpdateResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Handle Users - delete
func writeUsersDeleteResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeUsersDeleteResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// Handle writing person info response
func writePersonInfoResponse(rw http.ResponseWriter, req *http.Request) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := players.List()
	if err != nil {
		writeMessageResponse(rw, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	numberOfPeople := len(listOfPeople)

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(numberOfPeopleResponseJSON{
		NumberOfPeople: numberOfPeople,
	})
}

// Handle PreflightRequest
func writePreflightRequest(rw http.ResponseWriter, req *http.Request) {
	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}

// Handle writing the GET list of People response
func writeGetListOfPeopleResponse(rw http.ResponseWriter, req *http.Request) {

	setHeaders(rw, req)

	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := players.List()
	if err != nil {
		writeMessageResponse(rw, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listPeopleResponseJSON{
		People: listOfPeople,
	})
}

// Write the GET person response
func writeGetPersonDetailsResponse(rw http.ResponseWriter, req *http.Request, idString string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		writeMessageResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	person, err := players.Details(id)

	if err != nil {
		writeMessageResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(personDetailsResponseJSON{
		Person: *person,
	})
}

// Write the POST Add Person response
func writePostAddPersonResponse(rw http.ResponseWriter, req *http.Request) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	var p players.Person

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted in Add Person request"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &p)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse person data for Add Person request"))
		clientError++
		return
	}

	err = players.AddPerson(p)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create a new person. err:%s", err))
		serverError++
		return
	}

	setHeaders(rw, req)

	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Write the DELETE person response
func writeDeletePersonResponse(rw http.ResponseWriter, req *http.Request, idString string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	// Convert the ID into a number
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Logger.Printf(err.Error())
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("The ID:%s is not a number", idString))
		clientError++
		return
	}

	players.Delete(id)

	setHeaders(rw, req)

	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Write the GET metrics response
func writeGetMetricsResponse(rw http.ResponseWriter, req *http.Request) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid username and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(metricsResponseJSON{
		ClientSuccess:             clientSuccess,
		ClientError:               clientError,
		ClientAuthenticationError: clientAuthenticationError,
		ServerError:               serverError,
	})
}

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	// Authenticate
	r.HandleFunc(baseURL+"/authenticate", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeAuthenticateResponse(w, req)
		})).Methods(http.MethodPost)

	// Users - getAll
	r.HandleFunc(baseURL+"/users", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUsersGetAllResponse(w, req)
		})).Methods(http.MethodGet)

	// Users - getById
	r.HandleFunc(baseURL+"/users/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUsersGetByIDResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	// Users - create
	r.HandleFunc(baseURL+"/users", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUsersCreateResponse(w, req)
		})).Methods(http.MethodPost)

	// Users - update
	r.HandleFunc(baseURL+"/users", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUsersUpdateResponse(w, req)
		})).Methods(http.MethodPut)

	// Users - delete
	r.HandleFunc(baseURL+"/users", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUsersDeleteResponse(w, req)
		})).Methods(http.MethodDelete)

	// PersonInfo
	r.HandleFunc(baseURL+"/personinfo", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writePersonInfoResponse(w, req)
		})).Methods(http.MethodGet)

	// ListPeople
	r.HandleFunc(baseURL+"/people", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetListOfPeopleResponse(w, req)
		})).Methods(http.MethodGet)

	r.HandleFunc(baseURL+"/people", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writePreflightRequest(w, req)
		})).Methods(http.MethodOptions)

	// Person Details
	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetPersonDetailsResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	// Add Person
	r.HandleFunc(baseURL+"/person", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writePostAddPersonResponse(w, req)
		})).Methods(http.MethodPost)

	// Delete Person
	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeDeletePersonResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	// Metrics
	r.HandleFunc(baseURL+"/metrics", logger.LogHandler(
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

	if u != username {
		logger.Logger.Printf("checkUser: FAIL: username = %s, u = %s\n", username, u)
		return false
	} else if p != password {
		logger.Logger.Printf("checkUser: FAIL: password = %s %s\n", password, p)
		return false
	}

	logger.Logger.Printf("checkUser: OK\n")
	return true
}

func main() {

	logger.Logger.Printf("Players Server: 2018-01-31 13:30")
	var ok bool

	username, ok = os.LookupEnv("USERNAME")
	if !ok {
		username = "foo"
	}

	password, ok = os.LookupEnv("PASSWORD")
	if !ok {
		password = "bar"
	}

	baseURL, ok = os.LookupEnv("BASEURL")
	if !ok {
		baseURL = ""
	}

	portstring, ok := os.LookupEnv("PORT")
	if !ok {
		portstring = "4201"
	}
	port, err := strconv.Atoi(portstring)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Registering Router and setting Handlers")
	router := mux.NewRouter()
	setupHandlers(router)

	logger.Logger.Printf("Username = %s, Password = %s", username, password)
	logger.Logger.Printf("Listening to base URL: '%s' port: %d", baseURL, port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Success")
}
