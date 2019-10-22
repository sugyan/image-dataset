package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sugyan/image-dataset/web/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	app, err := app.NewApp(projectID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, app.Handler()); err != nil {
		log.Fatal(err)
	}
}
