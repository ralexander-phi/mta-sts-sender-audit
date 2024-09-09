package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LogLines struct {
	gorm.Model
	Uuid string
	Line string
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

	// Log received message
	result := db.Create(&LogLines{Uuid: userId, Line: "Message Received"})
	if result.Error != nil {
		panic(fmt.Sprintf("Unable to insert log: %v\n", result.Error))
	}
}
