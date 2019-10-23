package app

import (
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

func (app *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	results := []*imageResponse{}

	query := query(r)
	iter := app.dsClient.Run(r.Context(), query)
	for {
		var image entity.Image
		key, err := iter.Next(&image)
		if err != nil {
			if err == iterator.Done {
				break
			} else {
				log.Printf("failed to fetch data: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
		results = append(results, &imageResponse{
			ID:        key.Name,
			ImageURL:  image.ImageURL,
			Size:      image.Size,
			LabelName: image.LabelName,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&results); err != nil {
		log.Printf("failed to encode user info: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *App) userinfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client, err := app.firebase.Auth(ctx)
	if err != nil {
		log.Printf("failed to create auth client: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user, err := client.GetUser(ctx, app.uid(ctx))
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

func query(r *http.Request) *datastore.Query {
	q := r.URL.Query()
	query := datastore.NewQuery(entity.KindNameImage)
	if q.Get("key") != "" {
		query = query.Filter("__key__ >=", datastore.NameKey(entity.KindNameImage, q.Get("key"), nil))
	}
	return query.Limit(100)
}
