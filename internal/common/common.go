package common

import (
	"os"
	"os/user"
	"runtime"
	"unicode"

	codeerror "github.com/rsmaxwell/players-api/internal/codeerror"
)

// Reference type
type Reference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

var (
	// RootDir directory
	RootDir string

	// MetricsData containing metrics
	MetricsData Metrics
)

// Metrics structure
type Metrics struct {
	ClientSuccess             int `json:"clientSuccess"`
	ClientError               int `json:"clientError"`
	ClientAuthenticationError int `json:"clientAuthenticationError"`
	ServerError               int `json:"serverError"`
}

// HomeDir returns the home directory
func HomeDir() string {

	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}

	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func init() {

	home, ok := os.LookupEnv("PLAYERS_API_HOME")
	if !ok {
		home = HomeDir()
	}

	dirname, ok := os.LookupEnv("PLAYERS_API_DIRNAME")
	if !ok {
		dirname = "/players-api"
	}

	RootDir = home + dirname
}

// CheckCharactersInID checks the characters are valid for an ID
func CheckCharactersInID(s string) error {
	for _, r := range s {

		ok := false
		if unicode.IsLetter(r) {
			ok = true
		} else if unicode.IsDigit(r) {
			ok = true
		}

		if !ok {
			return codeerror.NewBadRequest("Invalid ID")
		}
	}
	return nil
}

// Contains function
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// EqualArrayOfStrings tells whether a and b contain the same elements NOT in-order order
func EqualArrayOfStrings(x, y []string) bool {

	if x == nil {
		if y == nil {
			return true
		}
		return false
	} else if y == nil {
		return false
	}

	if len(x) != len(y) {
		return false
	}

	xMap := make(map[string]int)
	yMap := make(map[string]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}

// EqualArrayOfStrings2 tells whether a and b contain the same elements IN order
func EqualArrayOfStrings2(x, y []string) bool {

	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

// SubtractLists Subtract the players on courts away from the list of players
func SubtractLists(listOfPlayers, players []string, text string) ([]string, error) {

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
