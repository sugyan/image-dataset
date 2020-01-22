package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gcp struct {
	csClient   *storage.Client
	fsClient   *firestore.Client
	bucketName string
}

func newGcp(projectID string) (*gcp, error) {
	ctx := context.Background()
	csClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucketName := projectID + ".appspot.com"
	if os.Getenv("DEVELOPMENT") != "" {
		bucketName = "staging." + bucketName
	}
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &gcp{
		csClient:   csClient,
		fsClient:   fsClient,
		bucketName: bucketName,
	}, nil
}

func (g *gcp) upload(filepath string) error {
	// load json file
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	data := &data{}
	if err := json.NewDecoder(jsonFile).Decode(data); err != nil {
		return err
	}
	// skip if image is too small
	if data.Size < 256 {
		return nil
	}

	// load image file
	name := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
	imageFile, err := os.Open(path.Join(path.Dir(filepath), name+".jpg"))
	if err != nil {
		return err
	}
	defer imageFile.Close()

	// calculate key name
	hash := md5.New()
	hash.Write([]byte(name))
	keyName := hex.EncodeToString(hash.Sum(nil))

	ctx := context.Background()
	if err := g.writeCS(ctx, keyName, imageFile); err != nil {
		return err
	}
	if err := g.writeFS(ctx, keyName, data); err != nil {
		return err
	}
	return nil
}

func (g *gcp) writeCS(ctx context.Context, objectName string, image io.Reader) error {
	obj := g.csClient.Bucket(g.bucketName).Object(fmt.Sprintf("images/%s", objectName))

	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, image); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

func (g *gcp) writeFS(ctx context.Context, keyName string, data *data) error {
	publishedAt, err := time.Parse("2006-01-02T15:04:05", data.Meta.PublishedAt)
	if err != nil {
		return err
	}
	parts := make([]int, 136)
	for i := 0; i < 68; i++ {
		parts[i*2], parts[i*2+1] = data.Parts[i][0], data.Parts[i][1]
	}
	metaData := map[string]interface{}{
		"angle": data.Angle,
	}
	if photoID, err := strconv.Atoi(data.Meta.PhotoID); err == nil {
		metaData["photo_id"] = photoID
	}
	if labelID, err := strconv.Atoi(data.Meta.LabelID); err == nil {
		metaData["label_id"] = labelID
	}
	meta, err := json.Marshal(&metaData)
	if err != nil {
		return err
	}

	return g.fsClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var image entity.Image
		docRef := g.fsClient.Collection(entity.KindNameImage).Doc(keyName)
		document, err := tx.Get(docRef)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				image = entity.Image{
					Status:    entity.StatusReady,
					CreatedAt: time.Now(),
				}
				// Update counts
				for i, size := range []int{256, 512, 1024} {
					if data.Size >= size {
						docID := []string{"0256", "0512", "1024"}[i]
						ref := g.fsClient.Collection(entity.KindNameCount).Doc(docID)
						if err := tx.Update(ref, []firestore.Update{
							{Path: entity.StatusReady.Path(), Value: firestore.Increment(1)},
						}); err != nil {
							return err
						}
					}
				}
			} else {
				return err
			}
		} else {
			if err := document.DataTo(&image); err != nil {
				return err
			}
		}
		image.ID = keyName
		image.ImageURL = fmt.Sprintf("https://storage.googleapis.com/%s/images/%s", g.bucketName, keyName)
		image.SourceURL = data.Meta.SourceURL
		image.PhotoURL = data.Meta.PhotoURL
		image.Size = data.Size
		image.Size0256 = data.Size >= 256
		image.Size0512 = data.Size >= 512
		image.Size1024 = data.Size >= 1024
		image.Parts = parts
		image.LabelName = data.Meta.LabelName
		image.PublishedAt = publishedAt
		image.UpdatedAt = time.Now()
		image.Meta = meta
		return tx.Set(docRef, &image)
	})
}
