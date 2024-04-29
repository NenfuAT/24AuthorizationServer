package model

type Session struct {
	Client              string
	State               string
	Scopes              string
	RedirectUri         string
	CodeChallenge       string
	CodeChallengeMethod string
	// OIDC用
	Nonce string
	// IDトークンを払い出すか否か、trueならIDトークンもfalseならOAuthでトークンだけ払い出す
	Oidc bool
}

var SessionList = make(map[string]Session)
