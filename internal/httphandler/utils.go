package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/common"
	"github.com/rsmaxwell/players-api/internal/debug"
	"github.com/rsmaxwell/players-api/internal/response"
)

// messageResponse structure
type messageResponse struct {
	Message string `json:"message"`
}

var (
	contextPath = "/players-api"

	pkg = debug.NewPackage("httphandler")

	functionMiddleware        = debug.NewFunction(pkg, "Middleware")
	functionCheckAccessToken  = debug.NewFunction(pkg, "checkAccessToken")
	functionCheckRefreshToken = debug.NewFunction(pkg, "checkRefreshToken")
	functionAuthenticate      = debug.NewFunction(pkg, "Authenticate")
)

// writeResponseMessage method
func writeResponseMessage(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string, message string) {
	writeResponse(r, req, statusCode, qualifier)
	json.NewEncoder(r).Encode(messageResponse{
		Message: message,
	})
}

// writeResponseObject method
func writeResponseObject(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string, object interface{}) {
	writeResponse(r, req, statusCode, qualifier)
	json.NewEncoder(r).Encode(object)
}

// writeResponse method
func writeResponse(r http.ResponseWriter, req *http.Request, statusCode int, qualifier string) {

	common.MetricsData.StatusCodes[statusCode]++

	if statusCode == http.StatusOK {

		origin := req.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:4200"
		}

		r.Header().Set("Content-Type", "application/json")
		r.Header().Set("Access-Control-Allow-Origin", origin)
		r.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		r.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Origin, Authorization")

	} else if statusCode == http.StatusUnauthorized {

		r.Header().Set("WWW-Authenticate", "Basic realm=\"players-api: "+qualifier+"\"")

	} else {

	}

	r.WriteHeader(statusCode)
}

// writeResponseError function
func writeResponseError(rw http.ResponseWriter, req *http.Request, err error) {
	if err != nil {
		if serr, ok := err.(*codeerror.CodeError); ok {
			writeResponseMessage(rw, req, serr.Code(), serr.Qualifier(), serr.Error())
			return
		}

		writeResponseMessage(rw, req, http.StatusInternalServerError, "", "error")
		return
	}
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}

// RefreshClaims struct
type RefreshClaims struct {
	UserID    string
	ExpiresAt int64
	Count     int
}

// setRefreshClaims function
func setRefreshClaims(token *jwt.Token, claims *RefreshClaims) {

	c := token.Claims.(jwt.MapClaims)

	c["sub"] = claims.UserID
	c["exp"] = claims.ExpiresAt
	c["Count"] = claims.Count
}

// getRefreshClaims function
func getRefreshClaims(claims jwt.MapClaims, refreshClaims *RefreshClaims) error {

	var ok bool

	value := claims["sub"]
	refreshClaims.UserID, ok = value.(string)
	if !ok {
		return fmt.Errorf("The 'sub' value is not a 'string': %v, %t", value, value)
	}

	value = claims["exp"]
	float64Value, ok := value.(float64)
	if !ok {
		return fmt.Errorf("The 'exp' value is not a 'int64': %v, %t", value, value)
	}
	refreshClaims.ExpiresAt = int64(float64Value)

	value = claims["Count"]
	float64Value, ok = value.(float64)
	if !ok {
		return fmt.Errorf("The 'Count' value is not a 'number': %v, %t", value, value)
	}
	refreshClaims.Count = int(float64Value)

	return nil
}

// AccessClaims struct
type AccessClaims struct {
	UserID    string
	ExpiresAt int64
	Role      string
	FirstName string
	LastName  string
}

// setAccessClaims function
func setAccessClaims(token *jwt.Token, claims *AccessClaims) {

	c := token.Claims.(jwt.MapClaims)

	c["sub"] = claims.UserID
	c["exp"] = claims.ExpiresAt
	c["Role"] = claims.Role
	c["FirstName"] = claims.FirstName
	c["LastName"] = claims.LastName
}

// getAccessClaims function
func getAccessClaims(claims jwt.MapClaims, accessClaims *AccessClaims) error {

	var ok bool

	value := claims["sub"]
	accessClaims.UserID, ok = value.(string)
	if !ok {
		return fmt.Errorf("The 'sub' value is not a 'string': %v, %t", value, value)
	}

	value = claims["exp"]
	float64Value, ok := value.(float64)
	if !ok {
		return fmt.Errorf("The 'exp' value is not a 'int64': %v, %t", value, value)
	}
	accessClaims.ExpiresAt = int64(float64Value)

	value = claims["Role"]
	accessClaims.Role, ok = value.(string)
	if !ok {
		return fmt.Errorf("The 'Role' value is not a 'string': %v, %t", value, value)
	}

	value = claims["FirstName"]
	accessClaims.FirstName, ok = value.(string)
	if !ok {
		return fmt.Errorf("The 'FirstName' value is not a 'string': %v, %t", value, value)
	}

	value = claims["LastName"]
	accessClaims.LastName, ok = value.(string)
	if !ok {
		return fmt.Errorf("The 'LastName' value is not a 'string': %v, %t", value, value)
	}

	return nil
}

// checkAccessToken method
func checkAccessToken(req *http.Request) (*AccessClaims, error) {
	f := functionCheckAccessToken

	// ********************************************************************
	// * Get the access token from the Header
	// ********************************************************************
	authorizationString := req.Header.Get("Authorization")
	if authorizationString == "" {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	accessTokenString, err := stripBearerPrefixFromTokenString(authorizationString)
	if err != nil {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	claims := jwt.MapClaims{}
	accessToken, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, codeerror.NewUnauthorized("Not Authorized")
		}
		return nil, codeerror.NewBadRequest(err.Error())
	}
	if !accessToken.Valid {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	accessClaims := AccessClaims{}
	err = getAccessClaims(claims, &accessClaims)
	if err != nil {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	if f.Level() >= debug.VerboseLevel {
		f.DebugVerbose("accessClaims:")
		f.DebugVerbose("    UserID:    %s", accessClaims.UserID)
		f.DebugVerbose("    ExpiresAt: %d", accessClaims.ExpiresAt)
		f.DebugVerbose("    Role:      %s", accessClaims.Role)
		f.DebugVerbose("    FirstName: %s", accessClaims.FirstName)
		f.DebugVerbose("    LastName:  %s", accessClaims.LastName)
	}

	return &accessClaims, nil
}

// checkRefreshToken method
func checkRefreshToken(req *http.Request) (*RefreshClaims, error) {
	f := functionCheckRefreshToken

	cookie, err := req.Cookie("players-api")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, codeerror.NewUnauthorized("Not Authorized")
		}
		return nil, codeerror.NewInternalServerError(err.Error())
	}

	refreshTokenString := cookie.Value

	claims := jwt.MapClaims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, codeerror.NewUnauthorized("Not Authorized")
		}
		return nil, codeerror.NewBadRequest(err.Error())
	}
	if !refreshToken.Valid {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	refreshClaims := RefreshClaims{}
	err = getRefreshClaims(claims, &refreshClaims)
	if err != nil {
		return nil, codeerror.NewUnauthorized("Not Authorized")
	}

	if f.Level() >= debug.VerboseLevel {
		f.DebugVerbose("refreshClaims:")
		f.DebugVerbose("    UserID:    %s", refreshClaims.UserID)
		f.DebugVerbose("    ExpiresAt: %d", refreshClaims.ExpiresAt)
		f.DebugVerbose("    Count:     %d", refreshClaims.Count)
	}

	return &refreshClaims, nil
}

// SetupHandlers Handlers for REST API routes
func SetupHandlers(r *mux.Router) {

	s := r.PathPrefix("/players-api").Subrouter()

	s.HandleFunc("/users/authenticate", Authenticate).Methods(http.MethodPost)
	s.HandleFunc("/users/register", Register).Methods(http.MethodPost)
	s.HandleFunc("/users/refresh", Refresh).Methods(http.MethodPost)
	s.HandleFunc("/users", ListPeople).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", DeletePerson).Methods(http.MethodDelete)
	s.HandleFunc("/users/logout", Logout).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", GetPerson).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", UpdatePerson).Methods(http.MethodPut)
	s.HandleFunc("/users/player/{id}", UpdatePersonPlayer).Methods(http.MethodPut)
	s.HandleFunc("/users/role/{id}", UpdatePersonRole).Methods(http.MethodPut)
	s.HandleFunc("/users/move", PostMove).Methods(http.MethodPost)

	s.HandleFunc("/court", ListCourts).Methods(http.MethodGet)
	s.HandleFunc("/court/{id}", GetCourt).Methods(http.MethodGet)
	s.HandleFunc("/court", CreateCourt).Methods(http.MethodPost)
	s.HandleFunc("/court/{id}", UpdateCourt).Methods(http.MethodPut)
	s.HandleFunc("/court/{id}", DeleteCourt).Methods(http.MethodDelete)

	s.HandleFunc("/queue", GetQueue).Methods(http.MethodGet)
	s.HandleFunc("/metrics", GetMetrics).Methods(http.MethodGet)

	r.NotFoundHandler = http.HandlerFunc(NotFound)
}

// Middleware method
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		f := functionMiddleware

		rw2 := response.New(rw)

		f.DebugRequest(req)
		h.ServeHTTP(rw2, req)
	})
}
