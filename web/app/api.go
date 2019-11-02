package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

const limit = 30

var (
	sizeMap = map[string]string{
		"256":  "Size0256",
		"512":  "Size0512",
		"1024": "Size1024",
	}
	sortMap = map[string]string{
		"id":        "__key__",
		"posted_at": "PostedAt",
		"name":      "LabelName",
	}
)

func (app *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	query, err := makeQuery(r)
	if err != nil {
		log.Printf("failed to make query: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
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

func makeQuery(r *http.Request) (*datastore.Query, error) {
	query := datastore.NewQuery(entity.KindNameImage).Limit(limit)
	if r.URL.Query().Get("size") != "" && r.URL.Query().Get("size") != "all" {
		if key, ok := sizeMap[r.URL.Query().Get("size")]; ok {
			query = query.Filter(fmt.Sprintf("%s =", key), true)
		} else {
			return nil, fmt.Errorf("invalid size query: %v", r.URL.Query().Get("size"))
		}
	}
	if r.URL.Query().Get("sort") != "" {
		if key, ok := sortMap[r.URL.Query().Get("sort")]; ok {
			if r.URL.Query().Get("order") == "desc" {
				key = "-" + key
			}
			query = query.Order(key)
		} else {
			return nil, fmt.Errorf("invalid sort query: %v", r.URL.Query().Get("sort"))
		}
	}
	if r.URL.Query().Get("limit") != "" {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			return nil, err
		}
		query = query.Limit(limit)
	}
	return query, nil
}
