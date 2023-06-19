package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/controllers"
	"github.com/toshnaik/CloudBoard/initializers"
	"github.com/toshnaik/CloudBoard/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.ConnectToRedis()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "website working"})
	})
	
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	
	r.POST("/logout", middleware.RequireAuth, controllers.Logout)
	r.POST("/refresh", middleware.RequireAuth, controllers.RefreshAccessToken)
	r.GET("/user", middleware.RequireAuth, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"user": ctx.MustGet("user_id")})
	})

	r.GET("/redis/test", func(ctx *gin.Context) {
		res, err := initializers.Redis.Incr(ctx, "counter").Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"counter": res})
	})

	r.Run(":" + os.Getenv("PORT"))
}
