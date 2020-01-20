package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/EeveesEyes/antilope-backend/db"
	"github.com/EeveesEyes/antilope-backend/models"
	"github.com/gin-gonic/gin"
	"log"
)

var IsProduction bool

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateUser(c *gin.Context) {
	user, err := GetValidUserFromRequest(c)
	if err != nil {
		if err.Error() == "weak password" {
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !db.UniqueEmail(user.Email) {
		c.JSON(400, gin.H{"error": "email is duplicate"})
		return
	}
	db.SaveUser(user)
	jwt, err := generateJWT(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	userAuthResponse := models.UserAuthResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    jwt,
	}
	c.JSON(201, gin.H{"user": userAuthResponse})
}

func GetAllUsers(c *gin.Context) {
	if IsProduction {
		c.Status(405)
		return
	}
	result, err := db.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"users": result})
	return
}

// Validators:
func GetValidUserFromRequest(c *gin.Context) (*models.User, error) {
	var userReq models.UserRequest
	bodyData, err := c.GetRawData()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyData, &userReq)
	if err != nil {
		return nil, err
	}

	if result, unmarshalErr := ValidatePassword(userReq.UserData.Password); unmarshalErr != nil {
		return nil, unmarshalErr
	} else if !result.Strong {
		c.JSON(400, gin.H{"error": result.Errors})
		return nil, fmt.Errorf("weak password")
	}

	hash, pepperID, err := HashPassword(userReq.UserData.Password)
	userReq.UserData.Password = GetRandString(len(userReq.UserData.Password))
	if err != nil {
		return nil, err
	}

	return models.NewUser(
		db.NextId,
		userReq.UserData.Username,
		userReq.UserData.Email,
		hash,
		pepperID), err
}
