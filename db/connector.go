package db

import (
	"fmt"
	"github.com/EeveesEyes/antilope-backend/models"
	"time"
)

var Users []*models.User
var JWTBlacklist = map[string]int64{}
var Secrets = map[int]models.Secret{}
var NextId = 1
var NextSecretId = 1

func GetUser(email string) (*models.User, error) {
	for _, v := range Users {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}
func GetUserById(userId int) (*models.User, error) {
	for _, v := range Users {
		if v.Id == userId {
			return v, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func SaveUser(user *models.User) {
	Users = append(Users, user)
	NextId++
}

func SaveSecret(secret models.Secret) {
	Secrets[secret.Id] = secret
	NextSecretId++
}

func UniqueEmail(email string) bool {
	for _, user := range Users {
		if user.Email == email {
			return false
		}
	}
	return true
}

func GetSecret(secretId int) (secret models.Secret) {
	return Secrets[secretId]
}

func AddJWTToBlacklist(tokenString string, expiration int64) {
	JWTBlacklist[tokenString] = expiration
	fmt.Println(JWTBlacklist)
}

// returns true if JWT is back-listed
func CheckJWTinBlacklist(tokenString string) (ok bool) {
	_, ok = JWTBlacklist[tokenString]
	return ok
}

func RemoveExpiredJWTsFromBlacklist() (deletedCount int) {
	for k, v := range JWTBlacklist {
		if v < time.Now().Unix() {
			delete(JWTBlacklist, k)
			deletedCount++
		}
	}
	return deletedCount
}

// Debug:

func GetAllUsers() ([]*models.User, error) {
	return Users, nil
}

func DeleteAllUsers() error {
	Users = []*models.User{}
	return nil
}
