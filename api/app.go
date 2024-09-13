package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogLines struct {
	gorm.Model
	Uuid    string
	Service string
	Line    string
}

type PollForm struct {
	Users string `form:"users"`
}

func serve() {

	MTA_STS_POLICY_DOC := `version: STSv1
mode: enforce
mx: *.audit.alexsci.com
max_age: 604800`

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
	r.GET("/.well-known/mta-sts.txt", func(c *gin.Context) {
		// Log info about the connecting client
		logLine := fmt.Sprintf(`HTTP connection:
			X-Forwarded-For: %s
			Host: %s
			User Agent: %s`,
			c.GetHeader("X-Forwarded-For"),
			c.Request.Host,
			c.GetHeader("User-Agent"))

		result := db.Create(&LogLines{
			Uuid:    "",
			Service: "HTTPS",
			Line:    logLine,
		})
		if result.Error != nil {
			c.JSON(500, gin.H{"message": "Unable to log policy request"})
			c.Abort()
			fmt.Println("Unable to get: ", result.Error)
			return
		}

		c.String(200, MTA_STS_POLICY_DOC)
		return
	})
	r.POST("/poll", func(c *gin.Context) {
		response := make(map[string][]map[string]string)
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
				response[line.Uuid] = []map[string]string{}
			}
			data := map[string]string{}
			data["service"] = line.Service
			data["when"] = line.CreatedAt.Format(time.RFC3339)
			data["line"] = line.Line
			response[line.Uuid] = append(response[line.Uuid], data)
		}
		c.JSON(200, response)
	})
	r.Run()
}

func healthcheck() {
	resp, err := http.Get("http://127.0.0.1:8080/health")
	if err != nil {
		fmt.Printf("Bad GET: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	// Just look for 200 OK status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Bad status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}
	fmt.Println("UP")
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "health" {
		healthcheck()
	} else {
		serve()
	}
}
