package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

type queryFilter struct {
	path  string
	op    string
	value interface{}
}

type queryOrder struct {
	Field string
	Desc  bool
}

type query struct {
	Filter []*queryFilter
	Order  *queryOrder
}

const limit = 30

var (
	sizeMap = map[string]string{
		"256":  "Size0256",
		"512":  "Size0512",
		"1024": "Size1024",
	}
	sortMap = map[string]string{
		"id":           "ID",
		"updated_at":   "UpdatedAt",
		"published_at": "PublishedAt",
	}
)

func (app *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	images, err := app.fetchImages(r)
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

func (app *App) updateImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data struct {
		Status int `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("failed to decode json: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err := app.updateImage(r.Context(), vars["id"], data.Status); err != nil {
		log.Printf("failed to update status: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) statsHandler(w http.ResponseWriter, r *http.Request) {
	query := app.fsClient.Collection(entity.KindNameCount).Query
	iter := query.Documents(r.Context())

	documents, err := iter.GetAll()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	log.Printf("%d documents", len(documents))
	for _, d := range documents {
		log.Printf("%v", d.Ref.ID)
		var count entity.Count
		if err := d.DataTo(&count); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Printf("%#v", count)
	}
	w.WriteHeader(http.StatusOK)
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

func (app *App) fetchImages(r *http.Request) ([]*imageResponse, error) {
	query, err := app.makeQuery(r)
	if err != nil {
		return nil, err
	}
	images := []*imageResponse{}
	iter := query.Documents(r.Context())
	for {
		document, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				return nil, err
			}
		}
		var image entity.Image
		if err := document.DataTo(&image); err != nil {
			return nil, err
		}
		images = append(images, &imageResponse{
			ID:          image.ID,
			ImageURL:    image.ImageURL,
			Size:        image.Size,
			Status:      int(image.Status),
			Parts:       image.Parts,
			LabelName:   image.LabelName,
			SourceURL:   image.SourceURL,
			PhotoURL:    image.PhotoURL,
			PublishedAt: image.PublishedAt.Unix(),
			UpdatedAt:   image.UpdatedAt.Unix(),
			Meta:        string(image.Meta),
		})
	}
	return images, nil
}

func (app *App) makeQuery(r *http.Request) (*firestore.Query, error) {
	values := r.URL.Query()
	collection := app.fsClient.Collection(entity.KindNameImage)
	query := collection.Query
	if values.Get("count") != "" {
		count, err := strconv.Atoi(values.Get("count"))
		if err != nil {
			return nil, err
		}
		query = query.Limit(count)
	} else {
		query = query.Limit(limit)
	}
	// `Where`
	{
		filters := []*queryFilter{}
		if values.Get("name") != "" {
			filters = append(filters, &queryFilter{
				path:  "LabelName",
				op:    "==",
				value: values.Get("name"),
			})
		}
		if values.Get("status") != "" && values.Get("status") != "all" {
			status, err := strconv.Atoi(values.Get("status"))
			if err != nil {
				return nil, err
			}
			filters = append(filters, &queryFilter{
				path:  "Status",
				op:    "==",
				value: status,
			})
		}
		if values.Get("size") != "" && values.Get("size") != "all" {
			if key, ok := sizeMap[values.Get("size")]; ok {
				filters = append(filters, &queryFilter{
					path:  key,
					op:    "==",
					value: true,
				})
			} else {
				return nil, fmt.Errorf("invalid size query: %v", values.Get("size"))
			}
		}
		for _, filter := range filters {
			query = query.Where(filter.path, filter.op, filter.value)
		}
	}
	// `Order`
	{
		if values.Get("sort") != "" {
			if path, ok := sortMap[values.Get("sort")]; ok {
				reverse := values.Get("reverse") == "true"
				if values.Get("order") == "desc" {
					reverse = !reverse
				}
				if values.Get("id") != "" {
					var op string
					if reverse {
						op = "<="
					} else {
						op = ">="
					}

					var image entity.Image
					document, err := collection.Doc(values.Get("id")).Get(context.Background())
					if err != nil {
						return nil, err
					}
					if err := document.DataTo(&image); err != nil {
						return nil, err
					}
					switch path {
					case "ID":
						query = query.Where(path, op, image.ID)
					case "PublishedAt":
						query = query.Where(path, op, image.PublishedAt)
					case "UpdatedAt":
						query = query.Where(path, op, image.UpdatedAt)
					}
				}
				if reverse {
					query = query.OrderBy(path, firestore.Desc)
				} else {
					query = query.OrderBy(path, firestore.Asc)
				}
			} else {
				return nil, fmt.Errorf("invalid sort query: %v", values.Get("sort"))
			}
		}
	}
	return &query, nil
}

func (app *App) updateImage(ctx context.Context, id string, status int) error {
	return app.fsClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := app.fsClient.Collection(entity.KindNameImage).Doc(id)
		doc, err := tx.Get(docRef)
		if err != nil {
			log.Printf("failed to get document: %s", err.Error())
			return err
		}
		var image entity.Image
		if err := doc.DataTo(&image); err != nil {
			log.Printf("failed to retrieve image from document: %s", err.Error())
			return err
		}
		if int(image.Status) != status {
			log.Printf("update status from %d to %d", image.Status, status)
			// Update counts
			for i, b := range []bool{image.Size0256, image.Size0512, image.Size1024} {
				if b {
					docID := []string{"0256", "0512", "1024"}[i]
					log.Printf("%v %v: %s", i, b, docID)
					ref := app.fsClient.Collection(entity.KindNameCount).Doc(docID)
					if err := tx.Update(ref, []firestore.Update{
						// TODO
						{Path: "Ready", Value: firestore.Increment(1)},
					}); err != nil {
						return err
					}
				}
			}
			// Update status
			image.Status = entity.Status(status)
			image.UpdatedAt = time.Now()
			if err := tx.Set(docRef, &image); err != nil {
				log.Printf("failed to set document: %s", err.Error())
				return err
			}

		}
		return nil
	})
}
