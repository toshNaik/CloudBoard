package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/toshnaik/CloudBoard/initializers"
)

// type clipboard struct {
// 	Content string `json:"data"`
// }

// func WriteCloudboard(c *gin.Context) {
// 	// get the user id from context
// 	userID := c.MustGet("user_id")

// 	var body clipboard
// 	if err := c.BindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// write the clipboard to redis for a day
// 	err := initializers.Redis.Set(context.TODO(), fmt.Sprint(userID), body.Content, time.Hour * 24).Err()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "cloudboard updated"})
// }

// func ReadCloudboard(c *gin.Context) {
// 	// get the user id from context
// 	userID := c.MustGet("user_id")

// 	// get the clipboard from redis
// 	clipboard, err := initializers.Redis.Get(context.Background(), fmt.Sprint(userID)).Result()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": clipboard})
// }

func WSCloudBoard(c *gin.Context) {
	userID := c.MustGet("user_id")

	// upgrade the connection to websocket
	conn, err := initializers.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer conn.Close()

	// read messages and publish on redis channel
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
			err = initializers.Redis.Publish(context.Background(), fmt.Sprint(userID), string(message)).Err()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
		}
	}()

	// subscribe to receive messages from redis channel and send them back to client
	for {
		sub := initializers.Redis.Subscribe(context.Background(), fmt.Sprint(userID))
		message, err := sub.ReceiveMessage(context.Background())
		if err != nil {
			fmt.Println("read:", err)
			break
		}

		byte_message := []byte(message.Payload)

		err = conn.WriteMessage(websocket.TextMessage, byte_message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}