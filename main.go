package main

import (
	"CookIt/models"
	"CookIt/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	models.ConnectDataBase()
	var router = make(chan *gin.Engine)
	go routes.GetRouter(router)
	var port string = os.Getenv("SERVER_PORT")
	server_addr := fmt.Sprintf(":%s", port)
	r := <-router
	r.Run(server_addr)
}
