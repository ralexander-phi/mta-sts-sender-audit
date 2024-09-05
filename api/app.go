package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
	"strings"
)

type PollForm struct {
	Users string `form:"users"`
}

func main() {
	// Valkey setup
	ctx := context.Background()
	db, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"valkey:6379"}})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"pong": "true",
		})
	})
	r.POST("/poll", func(c *gin.Context) {
		response := gin.H{}
		var form PollForm
		c.Bind(&form)
		userIds := strings.Split(form.Users, ",")
		if len(userIds) > 4 || len(userIds) == 0 {
			c.JSON(400, gin.H{"message": "Too many IDs"})
			c.Abort()
			return
		}
		results, err := db.Do(ctx, db.B().Mget().Key(userIds...).Build()).ToArray()
		if err != nil {
			c.JSON(500, gin.H{"message": "Unable to get"})
			c.Abort()
			fmt.Println("Unable to get: ", err)
			return
		}
		for i, result := range(results) {
			_, err := result.ToString()
			if valkey.IsValkeyNil(err) {
				// No value was set
				response[userIds[i]] = false
			} else {
				// A value was set
				response[userIds[i]] = true
			}
		}
		c.JSON(200, response)
	})
	r.Run()
}
