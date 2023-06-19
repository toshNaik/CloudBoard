package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toshnaik/CloudBoard/initializers"
)

type clipboard struct {
	Content string `json:"data"`
}

func WriteCloudboard(c *gin.Context) {
	// get the user id from context
	userID := c.MustGet("user_id")

	var body clipboard
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// write the clipboard to redis for a day
	err := initializers.Redis.Set(context.TODO(), fmt.Sprint(userID), body.Content, time.Hour * 24).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cloudboard updated"})
}

func ReadCloudboard(c *gin.Context) {
	// get the user id from context
	userID := c.MustGet("user_id")

	// get the clipboard from redis
	clipboard, err := initializers.Redis.Get(context.Background(), fmt.Sprint(userID)).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": clipboard})
}