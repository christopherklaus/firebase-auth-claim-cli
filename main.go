package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a valid user id, project and level. Use -h for help.")
		return
	}

	uid := flag.String("u", "", "UID to set claim for")
	set := flag.Bool("s", false, "Set claim for UID")
	projectId := flag.String("p", "", "Firebase project ID")
	userLevel := flag.String("l", "user", "User level to set (e.g. owner, paying, admin, user)")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: *projectId})

	if err != nil {
		log.Fatalf("Error initializing firebase: %v\n", err)
	}

	if *set {
		setClaimFor(*uid, *userLevel, app, ctx)
		return
	}

	getClaimFor(*uid, app, ctx)
}

func setClaimFor(uid string, level string, app *firebase.App, ctx context.Context) {
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v\n", err)
	}

	claims := map[string]interface{}{"level": level}
	err = client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		log.Fatalf("Error setting custom claims %v\n", err)
	}

	fmt.Printf("Setting new claim for %s: %s \n", uid, level)
}

func getClaimFor(uid string, app *firebase.App, ctx context.Context) {
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v\n", err)
	}

	user, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("Error getting user %v\n", err)
	}

	if level := user.CustomClaims["level"]; level != nil {
		log.Printf("Current claim for user %s: %s\n", uid, level)
		return
	}

	log.Printf("No claim found for user %s\n", uid)
}
