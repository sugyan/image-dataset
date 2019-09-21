package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type app struct {
	handler http.Handler
}

// NewApp function
func NewApp() http.Handler {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signin", signinHandler).Methods("POST")
	return &app{
		handler: router,
	}

}

func (app *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.handler.ServeHTTP(w, r)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := &struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("data: %v", data)
}
