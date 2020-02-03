package main

import (
	"flag"
	"github.com/EeveesEyes/antilope-backend/controllers"
	"github.com/EeveesEyes/antilope-backend/db"
	"github.com/EeveesEyes/antilope-backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

const IsProduction = false

func main() {
	password := flag.String("p", "", "the password used to decrypt/create the local storage file. Example: -p=password123")

	flag.Parse()

	if len(*password) <= 8 {
		panic("You chose a weak password!")
	}

	db.DecryptLocalSecrets(*password)
	router := gin.Default()
	api := router.Group("/api")
	routes.CreateRoutes(api)
	controllers.IsProduction = IsProduction

	handleGenericError(router.Run(":9000"))
	db.EncryptLocalSecrets(*password)

}

func handleGenericError(err error) {
	if err != nil {
		log.Println(err)
	}
}
