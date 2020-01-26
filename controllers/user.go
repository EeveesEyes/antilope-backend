package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/EeveesEyes/antilope-backend/db"
	"github.com/EeveesEyes/antilope-backend/models"
	"github.com/EeveesEyes/antilope-backend/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var IsProduction bool

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

/*
	expects json request body:
{
  "user": {
	"username": "stephan",
    "email": "test@test.de",
    "password": "Str0ngP4ssW%rd"
  }
}
	returns json response in format:
{
	"user": {
		"Username": "stephan",
		"Email": "test@test.de",
		"Token": "<JWT-Token>"
	}
}
*/
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
	jwt, err := util.GenerateJWT(user)
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

func DeleteAllUsers(c *gin.Context) {
	if IsProduction {
		c.Status(405)
		return
	}
	err := db.DeleteAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	c.Status(200)
	return
}

/*
	Expects json request body:
{
	"user": {
    "email": "test@test.de",
    "password": "Str0ngP4ssW%rd"
  }
}
	returns json response:
{
	"user": {
		"Username": "stephan",		// TODO: remove Username
		"Email": "test@test.de",	// TODO: remove email
		"Token": "<JWT-Token>"
	}
}
*/
func Login(c *gin.Context) {
	var userReq models.UserRequest
	bodyData, err := c.GetRawData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	err = json.Unmarshal(bodyData, &userReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	if userReq.UserData.Email == "" {
		c.JSON(400, gin.H{"error": "email required"})
		return
	} else if userReq.UserData.Password == "" {
		c.JSON(400, gin.H{"error": "password required"})
		return
	}
	validPW := false
	user, err := db.GetUser(userReq.UserData.Email)
	if err == nil {
		validPW = util.ValidatePassword(userReq.UserData.Password, *user)
	} else {
		// Timing can be used to get information about user presence, delay invalid requests
		emptyUser := models.User{}
		_ = util.ValidatePassword(userReq.UserData.Password, emptyUser)
	}
	if !validPW || err != nil {
		c.JSON(401, gin.H{"error": "invalid email or password"})
		return
	}
	jwt, err := util.GenerateJWT(user)
	userAuthResponse := models.UserAuthResponse{
		Username: user.Username, // TODO: remove
		Email:    user.Email,    // TODO: remove
		Token:    jwt,
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	c.JSON(201, gin.H{"user": userAuthResponse})
}

/*
	expects json request body:
{
	"token": "<JWT-Token>"
}
	returns empty body
*/
func Logout(c *gin.Context) {
	type LogoutReq struct {
		Token string `json:"token" binding:"required"`
	}
	var logoutReq LogoutReq
	bodyData, err := c.GetRawData()
	if err != nil {
		log.Println(err.Error())
		c.Status(500)
		return
	}
	err = json.Unmarshal(bodyData, &logoutReq)
	if err != nil {
		log.Println(err.Error())
		c.Status(400)
		return
	}
	if logoutReq.Token == "" {
		c.JSON(401, gin.H{"error": "token is required"})
		return
	}
	if err := util.InvalidateToken(logoutReq.Token); err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
	}

	c.Status(200)
}

func Test(c *gin.Context) {
	var body struct {
		Information string `json:"information" binding:"required"`
		Test        string `json:"test" binding:"required"`
		Test2       string `json:"test2" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(body.Information, body.Test, body.Test2)
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

	if result, unmarshalErr := util.TestPasswordStrength(userReq.UserData.Password); unmarshalErr != nil {
		return nil, unmarshalErr
	} else if !result.Strong {
		c.JSON(400, gin.H{"error": result.Errors})
		return nil, fmt.Errorf("weak password")
	}

	hash, pepperID, err := util.HashPassword(userReq.UserData.Password)
	userReq.UserData.Password = util.GetRandString(len(userReq.UserData.Password))
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
