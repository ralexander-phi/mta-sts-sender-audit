package main

import (
	"context"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"os"
	"strings"
)

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

	// Valkey setup
	ctx := context.Background()
	cache, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"valkey:6379"}})
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	// Insert key
	err = cache.Do(ctx, cache.B().Set().Key(userId).Value("1").Build()).Error()
	if err != nil {
		panic(err)
	}
}
