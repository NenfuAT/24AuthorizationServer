package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/NenfuAT/24AuthorizationServer/model"
	"github.com/lestrrat-go/jwx/jwk"
)

// 秘密鍵を読み込む
func readPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile("private-key.pem")
	if err != nil {
		return nil, err
	}
	keyblock, _ := pem.Decode(data)
	if keyblock == nil {
		return nil, fmt.Errorf("invalid private key data")
	}
	if keyblock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key type : %s", keyblock.Type)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(keyblock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// "ヘッダー.ペイロード"を作成する
func makeHeaderPayload() string {
	// ヘッダー
	var header = []byte(`{"alg":"RS256","kid": "12345678","typ":"JWT"}`)

	// ペイロード
	var payload = model.Payload{
		Iss:        "https://oreore.oidc.com",
		Azp:        model.ClientInfo.Id,
		Aud:        model.ClientInfo.Id,
		Sub:        model.TestUser.Sub,
		AtHash:     "PRzSZsEPQVqzY8xyB2ls5A",
		Nonce:      "abc",
		Name:       model.TestUser.NameJa,
		GivenName:  model.TestUser.GivenName,
		FamilyName: model.TestUser.FamilyName,
		Locale:     model.TestUser.Locale,
		Iat:        time.Now().Unix(),
		Exp:        time.Now().Unix() + model.ACCESS_TOKEN_DURATION,
	}
	payload_json, _ := json.Marshal(payload)

	//base64エンコード
	b64header := base64.RawURLEncoding.EncodeToString(header)
	b64payload := base64.RawURLEncoding.EncodeToString(payload_json)

	return fmt.Sprintf("%s.%s", b64header, b64payload)
}

// JWTを作成
func makeJWT() (string, error) {
	jwtString := makeHeaderPayload()

	privateKey, err := readPrivateKey()
	if err != nil {
		return "", err
	}
	err = privateKey.Validate()
	if err != nil {
		return "", fmt.Errorf("private key validate err : %s", err)
	}
	hasher := sha256.New()
	hasher.Write([]byte(jwtString))
	tokenHash := hasher.Sum(nil)

	// 署名を作成
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, tokenHash)
	if err != nil {
		return "", fmt.Errorf("sign by private key is err : %s", err)
	}
	enc := base64.RawURLEncoding.EncodeToString(signature)

	// "ヘッダー.ペイロード.署名"を作成して返す
	return fmt.Sprintf("%s.%s", jwtString, enc), nil
}

// 　JWKを作成してJSONにして返す
func makeJWK() []byte {

	data, _ := os.ReadFile("public-key.pem")
	keyset, _ := jwk.ParseKey(data, jwk.WithPEM(true))

	keyset.Set(jwk.KeyIDKey, "12345678")
	keyset.Set(jwk.AlgorithmKey, "RS256")
	keyset.Set(jwk.KeyUsageKey, "sig")

	jwk := map[string]interface{}{
		"keys": []interface{}{keyset},
	}
	buf, _ := json.MarshalIndent(jwk, "", "  ")
	return buf

}
