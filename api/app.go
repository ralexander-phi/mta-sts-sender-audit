package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

type LogLines struct {
	gorm.Model
	Uuid string
	Line string
}

type PollForm struct {
	Users string `form:"users"`
}

func main() {
	// Database setup
	dsn := fmt.Sprintf("host=db user=postgres password=%s dbname=audit port=5432 sslmode=disable TimeZone=UTC", os.Getenv("POSTGRES_PASSWORD"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	// Auto-migration
	db.AutoMigrate(&LogLines{})

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"pong": "true",
		})
	})
	r.POST("/poll", func(c *gin.Context) {
		response := make(map[string][]string)
		var form PollForm
		c.Bind(&form)
		userIds := strings.Split(form.Users, ",")
		if len(userIds) > 4 || len(userIds) == 0 {
			c.JSON(400, gin.H{"message": "Too many IDs"})
			c.Abort()
			return
		}

		var logLines []LogLines
		result := db.Where("uuid IN ?", userIds).Order("created_at, id").Find(&logLines)
		if result.Error != nil {
			c.JSON(500, gin.H{"message": "Unable to get"})
			c.Abort()
			fmt.Println("Unable to get: ", result.Error)
			return
		}

		for _, line := range logLines {
			if _, has := response[line.Uuid]; !has {
				response[line.Uuid] = []string{}
			}
			response[line.Uuid] = append(response[line.Uuid], line.Line)
		}
		c.JSON(200, response)
	})
	r.Run()
}
