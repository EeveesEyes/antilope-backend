package db

import "github.com/EeveesEyes/antilope-backend/models"

var Users []*models.User
var NextId = 0

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
