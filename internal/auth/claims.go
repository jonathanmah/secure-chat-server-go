package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetClaimFromToken(claimKey, token string, keyFunc func(*jwt.Token) (interface{}, error)) (string, error) {
	claims, err := ParseToken(token, keyFunc)
	if err != nil {
		return "", fmt.Errorf("Failed to extract claim %s from token", claimKey)
	}
	return claims[claimKey].(string), nil
}

// verifies and gets a claim value from raw activation token string
func GetClaimFromActivationToken(claimKey, token string) (string, error) {
	claimValue, err := GetClaimFromToken(claimKey, token, KeyFuncActivation)
	if err != nil {
		return "", err
	}
	return claimValue, nil
}

// verifies and gets a claim value from an access cookie
func GetClaimFromAccessCookie(claimKey string, r *http.Request) (string, error) {
	claims, err := ParseAccessCookie(r)
	if err != nil {
		return "", err
	}
	return claims[claimKey].(string), nil
}

// verifies and gets claim from query parameters
func GetClaimFromQueryParams(claimKey string, r *http.Request) (string, error) {
	claims, err := ParseQueryParams(r)
	if err != nil {
		return "", err
	}
	return claims[claimKey].(string), nil
}
