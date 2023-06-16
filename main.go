package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/controllers"
	"github.com/toshnaik/CloudBoard/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Run(":" + os.Getenv("PORT"))
}
