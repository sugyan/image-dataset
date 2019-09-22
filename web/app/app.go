package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
)

// App struct
type App struct {
	firebase *firebase.App
}

// NewApp function
func NewApp() (*App, error) {
	fbApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &App{
		firebase: fbApp,
	}, nil
}

// Handler method
func (app *App) Handler() http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signin", app.signinHandler).Methods("POST")

	return router
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := &struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Printf("failed to decode json: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	token, err := app.verifyIDToken(data.Token)
	if err != nil {
		log.Printf("failed to verify ID token: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	ok := false
	if admin, exist := token.Claims["admin"]; exist {
		if admin.(bool) {
			ok = true
		}
	}
	if !ok {
		log.Printf("user %s is not admin", token.UID)
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) verifyIDToken(token string) (*auth.Token, error) {
	client, err := app.firebase.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	return client.VerifyIDToken(context.Background(), token)
}
