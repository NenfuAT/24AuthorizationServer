package model

type JwkKey struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	Use string `json:"use"`
}
