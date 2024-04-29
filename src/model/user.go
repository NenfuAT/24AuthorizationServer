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

func InsertUser(u User) error {
	if err := db.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) User {
	var user User
	result := db.Where("email = ? ", email).First(&user)
	if result.Error != nil {
		return User{}
	}
	return user
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
