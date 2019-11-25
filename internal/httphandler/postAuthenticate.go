package httphandler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/model"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("<JWT_SECRET_KEY>")

// Claims is a struct that will be encoded to a JWT.
type Claims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

var (
	functionAuthenticate = debug.NewFunction(pkg, "Authenticate")
)

// Authenticate method
func Authenticate(rw http.ResponseWriter, req *http.Request) {
	f := functionAuthenticate

	id, password, _ := req.BasicAuth()

	f.DebugVerbose("id:       %s", id)
	f.DebugVerbose("password: %s", password)

	err := model.Authenticate(id, password)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	expirationTime := time.Now().Add(3 * time.Hour)
	claims := &Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token2",
		Value:   tokenString,
		Expires: expirationTime,
	})

	sess, err := globalSessions.SessionStart(rw, req)
	if err != nil {
		errorHandler(rw, req, err)
		return
	}
	defer sess.SessionRelease(rw)

	sess.Set("id", id)

	setHeaders(rw, req)
	rw.WriteHeader(http.StatusOK)
}
