package main

import (
	"context"
	"flag"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

var uid string

func init() {
	flag.StringVar(&uid, "uid", "", "target uid")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatal(err)
	}

	claims := map[string]interface{}{
		"admin": true,
	}
	if err := client.SetCustomUserClaims(ctx, uid, claims); err != nil {
		log.Fatal(err)
	}

	iter := client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v (%v): %v", user.UID, user.Email, user.CustomClaims)
	}
}
