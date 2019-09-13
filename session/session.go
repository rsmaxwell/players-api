package session

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Session structure
type Session struct {
	Username string    `json:"username"`
	Start    time.Time `json:"start"`
}

var (
	sessions map[string]*Session
)

func init() {
	sessions = map[string]*Session{}
}

// New creates a new session
func New(username string) (string, error) {
	s := new(Session)
	s.Username = username
	s.Start = time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	token := id.String()
	sessions[token] = s

	return token, nil
}
