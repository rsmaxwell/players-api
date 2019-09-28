package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rsmaxwell/players-api/codeError"
	"github.com/rsmaxwell/players-api/destination"
	"github.com/rsmaxwell/players-api/person"
)

var (
	clientSuccess             int
	clientError               int
	clientAuthenticationError int
	serverError               int
)

// messageResponse structure
type messageResponse struct {
	Message string `json:"message"`
}

// WriteResponse method
func WriteResponse(w http.ResponseWriter, httpStatus int, message string) {

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

// errorHandler function
func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		setHeaders(rw, req)
		if serr, ok := err.(*codeError.CodeError); ok {
			WriteResponse(rw, serr.Code(), serr.Error())
			clientError++
			return
		}

		WriteResponse(rw, http.StatusInternalServerError, "InternalServerError")
		clientError++
		return
	}
}

// Subtract the players on courts away from the list of players
func subtractLists(listOfPlayers, players []string, text string) ([]string, error) {

	l := []string{}
	for _, id := range listOfPlayers {

		found := false
		for _, id2 := range players {
			if id == id2 {
				found = true
				break
			}
		}

		if !found {
			l = append(l, id)
		}
	}

	return l, nil
}

// startup checks the state on disk is consistent
func startup() error {

	// Make a list of players
	listOfPeople, err := person.List()
	if err != nil {
		return err
	}

	listOfPlayers := []string{}
	for _, id := range listOfPeople {
		p, err := person.Load(id)
		if err != nil {
			return err
		}
		if p.Player {
			listOfPlayers = append(listOfPlayers, id)
		}
	}

	// Subtract the players on courts away from the list of players
	listOfCourts, err := destination.ListCourts()
	if err != nil {
		return err
	}
	for _, id := range listOfCourts {
		c, err := destination.LoadCourt(id)
		if err != nil {
			return err
		}

		text := fmt.Sprintf("Court[%s]", id)
		listOfPlayers, err = subtractLists(listOfPlayers, c.Container.Players, text)
		if err != nil {
			return err
		}
	}

	// Subtract the players waiting in the queue away from the list of players
	q, err := destination.LoadQueue()
	if err != nil {
		return err
	}

	listOfPlayers, err = subtractLists(listOfPlayers, q.Container.Players, "Queue")
	if err != nil {
		return err
	}

	// The list of players should now be empty, however add any remaining players to the waiting queue
	for _, id := range listOfPlayers {
		q.Container.Players = append(q.Container.Players, id)
	}

	// Save the updated queue
	err = q.Save()
	if err != nil {
		return err
	}

	return nil
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(r *mux.Router) {

	r.HandleFunc("/register",
		func(w http.ResponseWriter, req *http.Request) {
			Register(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/login",
		func(w http.ResponseWriter, req *http.Request) {
			Login(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			ListCourts(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			GetCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/court",
		func(w http.ResponseWriter, req *http.Request) {
			CreateCourt(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdateCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/court/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			DeleteCourt(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/person",
		func(w http.ResponseWriter, req *http.Request) {
			ListPeople(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			GetPerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodGet)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			UpdatePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodPut)

	r.HandleFunc("/person/{id}",
		func(w http.ResponseWriter, req *http.Request) {
			DeletePerson(w, req, mux.Vars(req)["id"])
		}).Methods(http.MethodDelete)

	r.HandleFunc("/metrics",
		func(w http.ResponseWriter, req *http.Request) {
			GetMetrics(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/queue",
		func(w http.ResponseWriter, req *http.Request) {
			GetQueue(w, req)
		}).Methods(http.MethodGet)

	r.HandleFunc("/move",
		func(w http.ResponseWriter, req *http.Request) {
			PostMove(w, req)
		}).Methods(http.MethodPost)

	r.HandleFunc("/queue",
		func(w http.ResponseWriter, req *http.Request) {
			GetQueue(w, req)
		}).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			NotFound(w, req)
		})
}
