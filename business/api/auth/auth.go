// Package auth provides authentication and authorization support.
package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zaouldyeck/webservice/foundation/logger"
)

// ErrForbidden is returned for auth issues.
var ErrForbidden = errors.New("action is not allowed")

// Claims reprents authorization via JWT.
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// HasRole checks if the specified role exists.
func (c Claims) HasRole(r string) bool {
	for _, role := range c.Roles {
		if role == r {
			return true
		}
	}
	return false
}

// KeyLookup declares a method set of behaviour for looking keys for JWT use.
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config is information we need to initialize auth.
type Config struct {
	Log       *logger.Logger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients.
type Auth struct {
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
}
