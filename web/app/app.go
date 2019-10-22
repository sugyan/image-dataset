package app

import (
	"context"
	"encoding/hex"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const sessionUser = "user"

// App struct
type App struct {
	firebase *firebase.App
	dsClient *datastore.Client
	session  sessions.Store
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
	dsClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &App{
		firebase: fbApp,
		dsClient: dsClient,
		session:  sessions.NewCookieStore(sessionKey),
	}, nil
}

// Handler method
func (app *App) Handler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/signin", app.signinHandler).Methods("POST")
	router.HandleFunc("/api/signout", app.signoutHandler).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/index", app.indexHandler).Methods("GET")
	api.HandleFunc("/userinfo", app.userinfoHandler).Methods("GET")
	api.Use(app.authMiddleware)

	return router
}
