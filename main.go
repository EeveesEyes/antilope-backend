package main

import (
	"github.com/EeveesEyes/antilope-backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	api := router.Group("/api")
	routes.CreateRoutes(api)

	handleGenericError(router.Run(":9000"))
}

func handleGenericError(err error) {
	if err != nil {
		log.Println(err)
	}
}
