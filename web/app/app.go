package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const sessionUser = "user"

// App struct
type App struct {
	firebase *firebase.App
	session  sessions.Store
}

// NewApp function
func NewApp() (*App, error) {
	fbApp, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	sessionKey, err := hex.DecodeString(os.Getenv("SESSION_KEY"))
	if err != nil {
		return nil, err
	}
	return &App{
		firebase: fbApp,
		session:  sessions.NewCookieStore(sessionKey),
	}, nil
}

// Handler method
func (app *App) Handler() http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signin", app.signinHandler).Methods("POST")
	api.HandleFunc("/user", app.userHandler).Methods("GET")

	return router
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// verify ID token
	data := &struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Printf("failed to decode json: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	token, err := app.verifyIDToken(r.Context(), data.Token)
	if err != nil {
		log.Printf("failed to verify ID token: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// authorize admin only
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
	// save to session
	session, err := app.session.Get(r, sessionUser)
	if err != nil {
		log.Printf("failed to get session: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	session.Values["uid"] = token.UID
	if err := app.session.Save(r, w, session); err != nil {
		log.Printf("failed to save session: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) userHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.session.Get(r, sessionUser)
	if err != nil {
		log.Printf("failed to get session: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	uid, exist := session.Values["uid"]
	if !exist {
		log.Printf("session uid does not exist")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	client, err := app.firebase.Auth(r.Context())
	if err != nil {
		log.Printf("failed to create auth client: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user, err := client.GetUser(r.Context(), uid.(string))
	if err != nil {
		log.Printf("failed to get user: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user.UserInfo); err != nil {
		log.Printf("failed to encode user info: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *App) verifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	client, err := app.firebase.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return client.VerifyIDToken(ctx, token)
}
