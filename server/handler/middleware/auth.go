package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/xerrors"
)

var (
	ErrUnauthorized = xerrors.New("unauthorized")
)

func FindToken(req *http.Request) (string, error) {
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		return "", ErrUnauthorized
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return "", ErrUnauthorized
	}

	return strings.TrimPrefix(tokenString, prefix), nil
}

type authenticatorClaims struct {
	UserName string
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) error {
	var claims authenticatorClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("my_secret_key"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}

	issuer, err := claims.GetIssuer()
	if err != nil {

	}

	if issuer

}
