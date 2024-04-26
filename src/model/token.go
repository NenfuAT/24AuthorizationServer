package model

type TokenCode struct {
	User      string
	ClientId  string
	Scopes    string
	ExpiresAt int64
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	IdToken     string `json:"id_token,omitempty"`
}

var TokenCodeList = make(map[string]TokenCode)
