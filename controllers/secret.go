package controllers

import (
	"github.com/EeveesEyes/antilope-backend/db"
	"github.com/EeveesEyes/antilope-backend/models"
	"github.com/EeveesEyes/antilope-backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSecret(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		return
	}
	var body struct {
		SecretId int `json:"secretId" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	secret := db.GetSecret(body.SecretId)
	if secret.AuthorizedUser == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "secretId not found"})
		return
	}
	user, err := db.GetUserById(secret.AuthorizedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if !util.ValidateJWT(authHeader, user.Email) {
		c.Status(401)
	}
	secretResponse := struct {
		SecretId          int
		SecretInformation string
	}{secret.Id, secret.Information}
	c.JSON(201, gin.H{"secret": secretResponse})
}

/*
	expects Authorization-Header
	expects json body:
	{
		"information": "<string>"
	}
*/
func CreateSecret(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		return
	}
	if !util.ValidateJWT(authHeader, "") {
		c.Status(401)
	}
	token, err := util.GetToken(authHeader)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	var body struct {
		Information string `json:"information" binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, _ := util.GetTokenClaims(token)
	user, err := db.GetUser(claims["iss"].(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	secret := models.Secret{
		Id:             db.NextSecretId,
		AuthorizedUser: user.Id,
		Information:    body.Information,
	}
	db.SaveSecret(secret)
	secretCreatedResponse := struct {
		SecretId int `json:"secretId"`
	}{SecretId: secret.Id}
	c.JSON(201, gin.H{"secret": secretCreatedResponse})
}
