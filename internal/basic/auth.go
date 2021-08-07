package basic

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	mySigningKey = []byte("<SESSION_SECRET_KEY>")
)

type MyJwtClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

// GenerateToken generates a jwt token
func GenerateToken(id int, expiresAt time.Duration) (string, error) {

	claims := MyJwtClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(expiresAt).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// ValidateToken validates the jwt token
func ValidateToken(signedToken string) (*MyJwtClaims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&MyJwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*MyJwtClaims)
	if !ok {
		err = errors.New("couldn't parse jwt claims")
		return nil, err
	}

	return claims, nil
}
