package app

import (
	"context"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const sessionUser = "user"

// App struct
type App struct {
	firebase   *firebase.App
	fsClient   *firestore.Client
	session    sessions.Store
	adminToken string
}

// NewApp function
func NewApp(projectID string) (*App, error) {
	ctx := context.Background()
	fbApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	sessionKey, err := hex.DecodeString(os.Getenv("SESSION_KEY"))
	if err != nil {
		return nil, err
	}
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &App{
		firebase:   fbApp,
		fsClient:   fsClient,
		session:    sessions.NewCookieStore(sessionKey),
		adminToken: os.Getenv("ADMIN_TOKEN"),
	}, nil
}

// Handler method
func (app *App) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/signin", app.signinHandler).Methods("POST")
	router.HandleFunc("/api/signout", app.signoutHandler).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/images", app.imagesHandler).Methods("GET")
	api.HandleFunc("/image/{id}", app.updateImageHandler).Methods("PUT")
	api.HandleFunc("/stats", app.statsHandler).Methods("GET")
	api.HandleFunc("/userinfo", app.userinfoHandler).Methods("GET")
	api.Use(app.authMiddleware)

	// wildcard endpoints
	router.PathPrefix("/").HandlerFunc(app.appHandler)

	return router
}

func (app *App) appHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate(w, "index.html"); err != nil {
		log.Printf("failed to render template: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func renderTemplate(w http.ResponseWriter, filename string) error {
	t, err := template.ParseFiles(filepath.Join("templates", filename))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html")
	return t.Execute(w, nil)
}
