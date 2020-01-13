package controllers

import (
	"encoding/json"
	"github.com/EeveesEyes/antilope-backend/models"
	"github.com/gin-gonic/gin"
	"log"
)

var Users []*models.User

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func CreateUser(c *gin.Context) {
	user, err := GetValidUserFromRequest(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	Users = append(Users, user)
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

	hash, pepperID, err := HashPassword(userReq.UserData.Password)
	if err != nil {
		return nil, err
	}

	return models.NewUser(
		userReq.UserData.Username,
		userReq.UserData.Password,
		hash,
		pepperID), err
}
