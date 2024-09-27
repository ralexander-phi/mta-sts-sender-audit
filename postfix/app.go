package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LogLines struct {
	gorm.Model
	Uuid    string
	Service string
	Line    string
	Public  bool
}

func main() {
	// Postfix input
	recipient := os.Getenv("ORIGINAL_RECIPIENT")
	userId, _, success := strings.Cut(recipient, "@")
	if !success {
		panic(fmt.Sprintf("Can't find user ID in: %s\n", recipient))
	}
	if len(userId) != 36 {
		panic(fmt.Sprintf("User ID looks invalid: %s\n", userId))
	}

	// Postgres password is saved to a file (as ENVs aren't passed down)
	content, err := ioutil.ReadFile("/postgres-password.txt")
	if err != nil {
		panic(fmt.Sprintf("Unable to get password from file: %v\n", err))
	}
	password := string(content)

	// Database setup
	dsn := fmt.Sprintf("host=db user=postgres password=%s dbname=audit port=5432 sslmode=disable TimeZone=UTC", password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database using %s: %v\n", password, err))
	}

	// Auto-migration
	db.AutoMigrate(&LogLines{})

	logLine := fmt.Sprintf(`Message Received:
		ENVID: %s
		SENDER: %s
		ORIGINAL_RECIPIENT: %s
		CLIENT_ADDRESS: %s
		CLIENT_HELO: %s
		CLIENT_HOSTNAME: %s
		CLIENT_PROTOCOL: %s`,
		os.Getenv("ENVID"),
		os.Getenv("SENDER"),
		os.Getenv("ORIGINAL_RECIPIENT"),
		os.Getenv("CLIENT_ADDRESS"),
		os.Getenv("CLIENT_HELO"),
		os.Getenv("CLIENT_HOSTNAME"),
		os.Getenv("CLIENT_PROTOCOL"))
	result := db.Create(&LogLines{
		Uuid:    userId,
		Service: "postfix",
		Line:    logLine,
		Public:  true,
	})
	if result.Error != nil {
		panic(fmt.Sprintf("Unable to insert log: %v\n", result.Error))
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result := db.Create(&LogLines{
			Uuid:    userId,
			Service: "postfix",
			Line:    fmt.Sprintf("MSG: %s", scanner.Text()),
			Public:  false,
		})
		if result.Error != nil {
			panic(fmt.Sprintf("Unable to insert log: %v\n", result.Error))
		}
	}

	if scanner.Err() != nil {
		panic(fmt.Sprintf("Unable to insert body: %v", scanner.Err()))
	}
}
