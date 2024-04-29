package model

type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	NameJa     string `json:"name_ja"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
}

func InsertUser(u User) error {
	if err := db.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) User {
	var user User
	if err := db.
		Where("email = ?", email).
		Find(&user).Error; err != nil {
		return User{}
	}
	return user
}

func GetUserByEmailAndPassword(email, password string) (User, error) {
	var user User
	result := db.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
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
