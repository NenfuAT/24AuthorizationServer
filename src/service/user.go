package service

import (
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/model"
)

func CreateUser(req model.User) (model.User, error) {

	result := model.GetUserByEmail(req.Email)

	if result.ID != "" {
		return model.User{}, nil
	} else {
		model.InsertUser(req)
	}
	return req, nil
}

func CheckEmail(email string) (int, map[string]string) {
	user := model.GetUserByEmail(email)
	if user.ID != "" {
		response := map[string]string{
			"Error": "Email address already in use",
		}
		return http.StatusBadRequest, response
	} else {
		return http.StatusOK, nil
	}
}
