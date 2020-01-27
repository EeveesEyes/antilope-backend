package routes

import (
	"github.com/EeveesEyes/antilope-backend/controllers"
	"github.com/gin-gonic/gin"
)

func CreateRoutes(router gin.IRouter) {
	router.GET("/ping", controllers.Ping)
	router.POST("/users/", controllers.CreateUser)
	router.GET("/users/", controllers.GetAllUsers)
	router.DELETE("/users/", controllers.DeleteAllUsers)
	router.POST("/users/login/", controllers.Login)
	router.POST("/users/logout/", controllers.Logout)
	router.GET("/secrets", controllers.GetSecret)
	router.POST("/secrets", controllers.CreateSecret)
}
