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

	app := app.NewApp()
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, app); err != nil {
		log.Fatal(err)
	}
}
