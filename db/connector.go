package db

import (
	"fmt"
	"github.com/EeveesEyes/antilope-backend/models"
)

var Users []*models.User
var NextId = 0

func GetUser(email string) (*models.User, error) {
	for _, v := range Users {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func SaveUser(user *models.User) {
	Users = append(Users, user)
	NextId++
}

func UniqueEmail(email string) bool {
	for _, user := range Users {
		if user.Email == email {
			return false
		}
	}
	return true
}

func GetAllUsers() ([]*models.User, error) {
	return Users, nil
}

func DeleteAllUsers() error {
	Users = []*models.User{}
	return nil
}
