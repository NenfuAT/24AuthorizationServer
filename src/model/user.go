package model

type User struct {
	Id         int
	Name       string
	Password   string
	Sub        string
	NameJa     string
	GivenName  string
	FamilyName string
	Locale     string
}

var TestUser = User{
	Id:         1111,
	Name:       "hoge",
	Password:   "password",
	Sub:        "11111111",
	NameJa:     "徳川慶喜",
	GivenName:  "慶喜",
	FamilyName: "徳川",
	Locale:     "ja",
}
