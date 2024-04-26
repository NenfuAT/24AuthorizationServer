package model

type Payload struct {
	Iss        string `json:"iss"`
	Azp        string `json:"azp"`
	Aud        string `json:"aud"`
	Sub        string `json:"sub"`
	AtHash     string `json:"at_hash"`
	Nonce      string `json:"nonce"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
	Iat        int64  `json:"iat"`
	Exp        int64  `json:"exp"`
}
