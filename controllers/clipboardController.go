package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

func WSCloudBoard(c *gin.Context) {
	userID := c.Param("user_id")
	conn, err := initializers.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add the connection to the map
	initializers.ChannelMutex.Lock()
	initializers.ChannelMap[fmt.Sprint(userID)] = conn
	initializers.ChannelMutex.Unlock()

	// read from the connection
	go func() {
		defer func() {
			// remove the connection from the map
			initializers.ChannelMutex.Lock()
			delete(initializers.ChannelMap, fmt.Sprint(userID))
			initializers.ChannelMutex.Unlock()
			conn.Close()
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	// publish clipboard data to the connection
	pubsub := initializers.Redis.Subscribe(context.Background(), fmt.Sprint(userID))
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		initializers.ChannelMutex.Lock()
		conn := initializers.ChannelMap[fmt.Sprint(userID)]
		initializers.ChannelMutex.Unlock()

		err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			return
		}
	}
}