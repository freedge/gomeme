package types

import "time"

// TokenFileToken is a token as found in the token file
type TokenFileToken struct {
	Token
	Created time.Time
}

// TokenFile contains all the acquired tokens
type TokenFile struct {
	Endpoint map[string]TokenFileToken
}
