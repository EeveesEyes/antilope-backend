package db

import "github.com/EeveesEyes/antilope-backend/models"

var Users []*models.User

func SaveUser(user *models.User) {
	Users = append(Users, user)
}

func UniqueEmail(email string) bool {
	for _, user := range Users {
		if user.Email == email {
			return false
		}
	}
	return true
}
