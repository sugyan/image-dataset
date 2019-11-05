package app

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

type queryFilter struct {
	Str   string
	Value interface{}
}

type queryOrder struct {
	Field string
	Desc  bool
}

type query struct {
	Filter []*queryFilter
	Order  *queryOrder
}

func init() {
	gob.Register(&query{})
}

func newQuery(r *http.Request) (*query, error) {
	query := &query{Filter: []*queryFilter{}}
	if r.URL.Query().Get("name") != "" {
		query.Filter = append(query.Filter, &queryFilter{
			Str:   "LabelName =",
			Value: r.URL.Query().Get("name"),
		})
	}
	if r.URL.Query().Get("size") != "" && r.URL.Query().Get("size") != "all" {
		if key, ok := sizeMap[r.URL.Query().Get("size")]; ok {
			query.Filter = append(query.Filter, &queryFilter{
				Str:   fmt.Sprintf("%s =", key),
				Value: true,
			})
		} else {
			return nil, fmt.Errorf("invalid size query: %v", r.URL.Query().Get("size"))
		}
	}
	if r.URL.Query().Get("sort") != "" {
		if key, ok := sortMap[r.URL.Query().Get("sort")]; ok {
			query.Order = &queryOrder{
				Field: key,
				Desc:  r.URL.Query().Get("order") == "desc",
			}
		} else {
			return nil, fmt.Errorf("invalid sort query: %v", r.URL.Query().Get("sort"))
		}
	}
	return query, nil
}

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
	}
)

func (app *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	query, err := func() (*datastore.Query, error) {
		id := r.URL.Query().Get("id")
		if id != "" {
			session, err := app.session.Get(r, sessionUser)
			if err != nil {
				return nil, err
			}
			q, ok := session.Values["query"].(*query)
			if !ok {
				q = &query{}
			}
			return app.makeQuery(q, r.URL.Query().Get("reverse") == "true", datastore.NameKey(entity.KindNameImage, id, nil))
		}
		query, err := newQuery(r)
		if err != nil {
			return nil, err
		}
		session, err := app.session.Get(r, sessionUser)
		if err != nil {
			return nil, err
		}
		session.Values["query"] = query
		if err := session.Save(r, w); err != nil {
			return nil, err
		}
		return app.makeQuery(query, false, nil)
	}()
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

func (app *App) makeQuery(q *query, reverse bool, key *datastore.Key) (*datastore.Query, error) {
	query := datastore.NewQuery(entity.KindNameImage).Limit(limit)
	if q.Filter != nil {
		for _, filter := range q.Filter {
			query = query.Filter(filter.Str, filter.Value)
		}
	}
	if q.Order != nil {
		field := q.Order.Field
		if q.Order.Desc {
			reverse = !reverse
		}
		if key != nil {
			var inequality string
			if reverse {
				inequality = "<="
			} else {
				inequality = ">="
			}
			if field != "__key__" {
				image := &entity.Image{}
				if err := app.dsClient.Get(context.Background(), key, image); err != nil {
					return nil, err
				}
				switch field {
				case "PostedAt":
					query = query.Filter(fmt.Sprintf("%s %s", field, inequality), image.PostedAt)
				}
			} else {
				query = query.Filter(fmt.Sprintf("__key__ %s", inequality), key)
			}
		}
		if reverse {
			field = "-" + field
		}
		query = query.Order(field)
	}
	return query, nil
}
