package model

type Client struct {
	Id          string
	Name        string
	RedirectURL string
	Secret      string
}

// クライアント情報をハードコード
var ClientInfo = Client{
	Id:          "1234",
	Name:        "test",
	RedirectURL: "http://localhost:8084/callback",
	Secret:      "secret",
}
