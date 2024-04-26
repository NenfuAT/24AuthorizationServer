package model

type AuthCode struct {
	User        string
	ClientId    string
	Scopes      string
	RedirectUri string
	ExpiresAt   int64
}

var AuthCodeList = make(map[string]AuthCode)
