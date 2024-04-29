package model

type User struct {
	ID         string
	Email      string
	Password   string
	NameJa     string
	GivenName  string
	FamilyName string
	Locale     string
}

var TestUser = User{
	ID:         "11111111",
	Email:      "hoge@gmail.com",
	Password:   "password",
	NameJa:     "徳川慶喜",
	GivenName:  "慶喜",
	FamilyName: "徳川",
	Locale:     "ja",
}
