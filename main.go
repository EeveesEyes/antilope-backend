package main

import (
	"github.com/EeveesEyes/antilope-backend/controllers"
	"github.com/EeveesEyes/antilope-backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

const IsProduction = false

func main() {
	router := gin.Default()
	api := router.Group("/api")
	routes.CreateRoutes(api)
	controllers.IsProduction = IsProduction

	handleGenericError(router.Run(":9000"))
}

func handleGenericError(err error) {
	if err != nil {
		log.Println(err)
	}
}
