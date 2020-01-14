package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type contextKey string

const (
	contextKeyUID contextKey = "uid"
	bearerPrefix  string     = "Bearer "
)

func (app *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, bearerPrefix) {
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			if token == app.adminToken {
				r = r.WithContext(context.WithValue(r.Context(), contextKeyUID, "admin"))
				next.ServeHTTP(w, r)
				return
			}
		}
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
		r = r.WithContext(context.WithValue(r.Context(), contextKeyUID, uid))
		next.ServeHTTP(w, r)
	})
}

func (app *App) uid(ctx context.Context) string {
	return ctx.Value(contextKeyUID).(string)
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// verify ID token
	var data struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
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

func (app *App) verifyIDToken(ctx context.Context, token string) (*auth.Token, error) {
	client, err := app.firebase.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return client.VerifyIDToken(ctx, token)
}

func (app *App) signoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.session.Get(r, sessionUser)
	if err != nil {
		log.Printf("failed to get session: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	if err := app.session.Save(r, w, session); err != nil {
		log.Printf("failed to save session: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
