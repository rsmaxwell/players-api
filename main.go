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
	"golang.org/x/crypto/bcrypt"

	"github.com/rsmaxwell/players-api/logger"
	"github.com/rsmaxwell/players-api/players"
	"github.com/rsmaxwell/players-api/session"
)

var (
	port                      int
	baseURL                   string
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// Authenticate Response
type authenticateResponse struct {
	Token string `json:"token"`
}

// Error Response
type messageResponse struct {
	Message string `json:"message"`
}

// metrics Response
type metricsResponse struct {
	ClientSuccess             int `json:"clientSuccess"`
	ClientError               int `json:"clientError"`
	ClientAuthenticationError int `json:"clientAuthenticationError"`
	ServerError               int `json:"serverError"`
}

// Court details Response
type courtDetailsResponse struct {
	Court players.Court `json:"court"`
}

// List courts Response
type listCourtsResponse struct {
	Courts []int `json:"courts"`
}

// People details Response
type personDetailsResponse struct {
	Person players.Person `json:"person"`
}

// List people Response
type listPeopleResponse struct {
	People []int `json:"people"`
}

// Handle writing error messager response.
func writeMessageResponse(w http.ResponseWriter, httpStatus int, message string) {
	logger.Logger.Printf("Response: %d %s - %s", httpStatus, http.StatusText(httpStatus), message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	json.NewEncoder(w).Encode(messageResponse{
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

// writeRegisterResponse
func writeRegisterResponse(rw http.ResponseWriter, req *http.Request) {

	var r players.RegisterRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse Person data"))
		clientError++
		return
	}

	err = players.RegisterPerson(r)
	if err != nil {
		logger.Logger.Printf("writeRegisterResponse: %s", err)
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("%s", err))
		serverError++
		return
	}

	setHeaders(rw, req)
	writeMessageResponse(rw, http.StatusOK, "ok")
}

// writeLoginResponse
func writeLoginResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeLoginResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	logger.Logger.Printf("writeLoginResponse(0): user:%s, password:%s", user, pass)

	if !checkUser(user, pass) {

		logger.Logger.Printf("writeLoginResponse(1)")

		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	logger.Logger.Printf("writeLoginResponse(2)")

	token, err := session.New(user)
	if err != nil {
		writeMessageResponse(rw, http.StatusInternalServerError, "Error creating session")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(authenticateResponse{
		Token: token,
	})
}

// Handle Court - getAll
func writeGetListOfCourtsResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeGetListOfCourtsResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfCourts, err := players.ListAllCourts()
	if err != nil {
		writeMessageResponse(rw, http.StatusInternalServerError, "Error getting list of courts")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listCourtsResponse{
		Courts: listOfCourts,
	})
}

// Handle Court - get by Id
func writeGetCourtByIDResponse(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeGetCourtByIDResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
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

	court, err := players.GetCourtDetails(id)

	if err != nil {
		writeMessageResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(courtDetailsResponse{
		Court: *court,
	})

	rw.WriteHeader(http.StatusOK)
}

// Handle Court - create
func writeCreateCourtResponse(rw http.ResponseWriter, req *http.Request) {

	logger.Logger.Printf("writeCreateCourtResponse")

	var r players.CreateCourtRequest

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted in Add Court request"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse data for CreateCourtRequest"))
		clientError++
		return
	}

	ok := session.CheckToken(r.Token)
	if !ok {
		writeMessageResponse(rw, http.StatusBadRequest, "Not Authorized")
		return
	}

	err = players.AddCourt(r.Court)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not create a new court. err:%s", err))
		serverError++
		return
	}

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}

// Handle Court - update
func writeUpdateCourtResponse(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeUpdateCourtResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
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

	var c players.JSONCourt

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse court data for id:%s", idString))
		clientError++
		return
	}

	_, err = players.UpdateCourt(id, c)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not update court:%s", idString))
		clientError++
		return
	}

	setHeaders(rw, req)
	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Handle Court - delete
func writeDeleteCourtResponse(rw http.ResponseWriter, req *http.Request, idString string) {

	logger.Logger.Printf("writeDeleteCourtResponse")

	// Check the user calling the service
	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
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

	err = players.DeleteCourt(id)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not delete court:%s", idString))
		clientError++
		return
	}

	setHeaders(rw, req)
	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Handle writing the GET list of People response
func writeGetListOfPeopleResponse(rw http.ResponseWriter, req *http.Request) {

	setHeaders(rw, req)

	user, pass, _ := req.BasicAuth()

	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	listOfPeople, err := players.ListAllPeople()
	if err != nil {
		writeMessageResponse(rw, http.StatusInternalServerError, "Error getting list of people")
		serverError++
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(listPeopleResponse{
		People: listOfPeople,
	})
}

// Write the GET person response
func writeGetPersonByIDResponse(rw http.ResponseWriter, req *http.Request, id string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	person, err := players.GetPersonDetails(id)
	if err != nil {
		writeMessageResponse(rw, http.StatusNotFound, "Not Found")
		clientError++
		return
	}

	setHeaders(rw, req)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(personDetailsResponse{
		Person: *person,
	})
}

// Update person
func writeUpdatePersonResponse(rw http.ResponseWriter, req *http.Request, id string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	var p players.JSONPerson

	limitedReader := &io.LimitedReader{R: req.Body, N: 20 * 1024}
	b, err := ioutil.ReadAll(limitedReader)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Too much data posted"))
		clientError++
		return
	}

	err = json.Unmarshal(b, &p)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not parse person data for person:%s", id))
		clientError++
		return
	}

	_, err = players.UpdatePerson(id, p)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not update person:%s", id))
		clientError++
		return
	}

	setHeaders(rw, req)
	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Write the DELETE person response
func writeDeletePersonResponse(rw http.ResponseWriter, req *http.Request, id string) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
		clientError++
		clientAuthenticationError++
		return
	}

	err := players.DeletePerson(id)
	if err != nil {
		writeMessageResponse(rw, http.StatusBadRequest, fmt.Sprintf("Could not delete person:%s", id))
		clientError++
		return
	}

	setHeaders(rw, req)
	writeMessageResponse(rw, http.StatusOK, "ok")
}

// Write the GET metrics response
func writeGetMetricsResponse(rw http.ResponseWriter, req *http.Request) {
	// Check the user calling the service
	user, pass, _ := req.BasicAuth()
	if !checkUser(user, pass) {
		writeMessageResponse(rw, http.StatusUnauthorized, "Invalid userID and/or password")
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

// Handlers for REST API routes
func setupHandlers(r *mux.Router) {

	// Register
	r.HandleFunc(baseURL+"/register", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeRegisterResponse(w, req)
		})).Methods(http.MethodPost)

	// Login
	r.HandleFunc(baseURL+"/login", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeLoginResponse(w, req)
		})).Methods(http.MethodGet)

	// -----[ Court ]---------------------------------------------

	// Court - getAll
	r.HandleFunc(baseURL+"/court", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetListOfCourtsResponse(w, req)
		})).Methods(http.MethodGet)

	// Court - getById
	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetCourtByIDResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	// Court - create
	r.HandleFunc(baseURL+"/court", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeCreateCourtResponse(w, req)
		})).Methods(http.MethodPost)

	// Court - update
	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUpdateCourtResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodPut)

	// Court - delete
	r.HandleFunc(baseURL+"/court/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeDeleteCourtResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	// -----[ People ]---------------------------------------------

	// People - GetAll
	r.HandleFunc(baseURL+"/person", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetListOfPeopleResponse(w, req)
		})).Methods(http.MethodGet)

	// Person - getById
	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeGetPersonByIDResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodGet)

	// Person - update
	r.HandleFunc(baseURL+"/person", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeUpdatePersonResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodPost)

	// Person - delete
	r.HandleFunc(baseURL+"/person/{id}", logger.LogHandler(
		func(w http.ResponseWriter, req *http.Request) {
			writeDeletePersonResponse(w, req, mux.Vars(req)["id"])
		})).Methods(http.MethodDelete)

	// -----[  ]---------------------------------------------

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
func checkUser(userID, password string) bool {

	person, err := players.GetPersonDetails(userID)
	if err != nil {
		logger.Logger.Printf("checkUser(0): userID or password do not match")
		return false
	}

	err = bcrypt.CompareHashAndPassword(person.HashedPassword, []byte(password))
	if err != nil {
		logger.Logger.Printf("checkUser(1): userID or password do not match")
		return false
	}

	logger.Logger.Printf("checkUser: OK\n")
	return true
}

func main() {

	logger.Logger.Printf("Players Server: 2018-01-31 13:30")
	var ok bool

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

	logger.Logger.Printf("Listening on port: %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Printf("Success")
}
