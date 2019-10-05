package session

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Session structure
type Session struct {
	UserID string    `json:"userID"`
	Start  time.Time `json:"start"`
}

var (
	sessions map[string]*Session
)

func init() {
	sessions = map[string]*Session{}
}

// New creates a new session
func New(userID string) (string, error) {

	s := new(Session)
	s.UserID = userID
	s.Start = time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	token := id.String()
	sessions[token] = s

	return token, nil
}

// LookupToken function
func LookupToken(token string) *Session {

	if s, ok := sessions[token]; ok {
		return s
	}

	return nil
}
