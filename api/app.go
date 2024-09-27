package main

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
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
	Public  bool
}

type PollForm struct {
	Users  string `form:"users"`
	Secret string `form:"secret"`
}

func serve() {
	const ONE_MB_IN_BYTES int64 = 1048576
	TLSRPT_GZIP := "application/tlsrpt+gzip"
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

	// Admin secret
	adminSecret := os.Getenv("ADMIN_SECRET")

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
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
			Public:  true,
		})
		if result.Error != nil {
			c.JSON(500, gin.H{"message": "Unable to log"})
			c.Abort()
			fmt.Println("Unable to get: ", result.Error)
			return
		}

		c.String(http.StatusOK, MTA_STS_POLICY_DOC)
		return
	})
	r.POST("/poll", func(c *gin.Context) {
		response := make(map[string][]map[string]string)
		var form PollForm
		c.Bind(&form)
		userIds := strings.Split(form.Users, ",")
		if len(userIds) > 4 || len(userIds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Too many IDs"})
			c.Abort()
			return
		}

		// Special admin access
		isAdmin := false
		if form.Secret != "" {
			actual := fmt.Sprintf("%x", sha256.Sum256([]byte(form.Secret)))
			if actual == adminSecret {
				isAdmin = true
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong secret"})
				c.Abort()
				return
			}
		}

		var logLines []LogLines
		query := db.Where("uuid IN ?", userIds)
		if !isAdmin {
			query = query.Where("public = ?", true)
		}
		result := query.Order("created_at, id").Find(&logLines)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get"})
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
		c.JSON(http.StatusOK, response)
	})
	r.POST("/tlsrpt", func(c *gin.Context) {
		// Protect against overly large reports
		reader := io.LimitReader(c.Request.Body, ONE_MB_IN_BYTES)

		if c.GetHeader("Content-Type") == TLSRPT_GZIP {
			// Need to decompress
			var err error
			reader, err = gzip.NewReader(reader)
			if err != nil {
				c.String(http.StatusBadRequest, "GZIP unreadable")
				c.Abort()
				return
			}
			// Project against overly large compressed inputs
			reader = io.LimitReader(reader, ONE_MB_IN_BYTES)
		}

		// Read the whole report
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			c.String(http.StatusBadRequest, "Report unreadable")
			c.Abort()
			return
		}

		result := db.Create(&LogLines{
			Uuid:    "",
			Service: "TLSRPT",
			Line:    string(body),
			Public:  true,
		})
		if result.Error != nil {
			c.String(500, "Unable to log")
			c.Abort()
			fmt.Println("Unable to get: ", result.Error)
			return
		}
		c.String(http.StatusOK, "OK")
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
