package service

import (
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
