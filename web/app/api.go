package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
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
	query, err := app.makeQuery(r)
	if err != nil {
		log.Printf("failed to make query: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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

func (app *App) fetchImages(ctx context.Context, query *firestore.Query) ([]*imageResponse, error) {
	images := []*imageResponse{}
	iter := query.Documents(ctx)
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
	query := collection.Limit(limit)
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
