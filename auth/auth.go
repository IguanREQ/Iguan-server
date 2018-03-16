package auth

import (
	"errors"
)

const (
	AuthTypeNone  byte = 1 << iota
	AuthTypeToken
	AuthTypeName
)

var (
	ErrParseAuthData = errors.New("Error while parsing auth data")
	ErrAccessDenied  = errors.New("Access denied")
)

// HTTP AuthCredentials headers
// https://tools.ietf.org/html/rfc6648#section-2
// Iguan-AuthType
// Iguan-Login
// Iguan-Password

type AuthCredentials struct {
	authType byte
	login    []byte
	password []byte
}

func (a *AuthCredentials) Valid() bool {
	if (a.authType & AuthTypeNone) != 0 {
		return true
	}

	if (a.authType&AuthTypeToken) != 0 && a.login == nil {
		return false
	}

	if (a.authType&AuthTypeName) != 0 && a.password == nil {
		return false
	}

	// TODO: аутентификация по реквизитам

	return true
}
