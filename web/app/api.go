package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

const limit = 30

func (app *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	query := datastore.NewQuery(entity.KindNameImage).Limit(limit)
	id := r.URL.Query().Get("id")
	if id != "" {
		idKey := datastore.NameKey(entity.KindNameImage, id, nil)
		if r.URL.Query().Get("order") == "desc" {
			query = query.Filter("__key__ <=", idKey).Order("-__key__")
		} else {
			query = query.Filter("__key__ >=", idKey).Order("__key__")
		}
	}
	// fetch forward and backward images
	images, err := app.fetchImages(r.Context(), query)
	if err != nil {
		log.Printf("failed to fetch data: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&images); err != nil {
		log.Printf("failed to encode user info: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *App) searchHandler(w http.ResponseWriter, r *http.Request) {
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
			Parts:     image.Parts,
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

func (app *App) fetchImages(ctx context.Context, query *datastore.Query) ([]*imageResponse, error) {
	images := []*imageResponse{}
	iter := app.dsClient.Run(ctx, query)
	for {
		var image entity.Image
		key, err := iter.Next(&image)
		if err != nil {
			if err == iterator.Done {
				break
			} else {
				return nil, err
			}
		}
		images = append(images, &imageResponse{
			ID:        key.Name,
			ImageURL:  image.ImageURL,
			Size:      image.Size,
			Parts:     image.Parts,
			LabelName: image.LabelName,
			SourceURL: image.SourceURL,
			PhotoURL:  image.PhotoURL,
			PostedAt:  image.PostedAt.Unix(),
			Meta:      string(image.Meta),
		})
	}
	return images, nil
}

func query(r *http.Request) *datastore.Query {
	query := datastore.NewQuery(entity.KindNameImage)
	// TODO: search query
	return query.Limit(100)
}
