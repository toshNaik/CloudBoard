package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/initializers"
	"github.com/toshnaik/CloudBoard/models"
	"github.com/toshnaik/CloudBoard/utils"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(c *gin.Context) {
	// get the token from the request header
	accessToken := c.Request.Header.Get("Authorization")

	// retrieve the token from the cookie if it is not present in the header
	if strings.HasPrefix(accessToken, "Bearer") {
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	} else {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no access token found"})
			return
		}
		accessToken = cookie
	}

	// validate the token
	claims, err := utils.VerifyToken(accessToken, os.Getenv("ACCESS_TOKEN_PUBLIC_KEY"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token",
																"message": err.Error()})
		return
	}
	// retrieve the user id from the Redis database
	userID, err := initializers.Redis.Get(context.TODO(), claims.TokenUuid).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired",
																"message": err.Error()})
		return
	}
	
	// check if the user still exists in the database
	var user models.User
	err = initializers.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found",
																"message": err.Error()})
		return
	}

	// set the user id and token uuid in the context
	c.Set("user_id", claims.UserID)
	c.Set("access_token_uuid", claims.TokenUuid)
	c.Next()
}