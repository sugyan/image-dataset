package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
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

	// load image file
	name := strings.TrimSuffix(path.Base(filepath), path.Ext(filepath))
	imageFile, err := os.Open(path.Join(path.Dir(filepath), name+".png"))
	if err != nil {
		return err
	}
	defer imageFile.Close()

	image, err := png.Decode(imageFile)
	if err != nil {
		return err
	}

	// calculate key name
	hash := md5.New()
	hash.Write([]byte(name))
	keyName := hex.EncodeToString(hash.Sum(nil))

	ctx := context.Background()
	if err := g.writeCS(ctx, keyName, image); err != nil {
		return err
	}
	if err := g.writeDS(ctx, keyName, data); err != nil {
		return err
	}
	return nil
}

func (g *gcp) writeCS(ctx context.Context, objectName string, image image.Image) error {
	obj := g.csClient.Bucket(g.bucketName).Object(fmt.Sprintf("images/%s", objectName))

	w := obj.NewWriter(ctx)
	if err := jpeg.Encode(w, image, &jpeg.Options{Quality: 100}); err != nil {
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

func (g *gcp) writeDS(ctx context.Context, keyName string, data *data) error {
	publishedAt, err := time.Parse("2006-01-02 15:04:05", data.Meta.PublishedAt)
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
	if faceID, err := strconv.Atoi(data.Meta.FaceID); err == nil {
		metaData["face_id"] = faceID
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

	var image entity.Image
	docRef := g.fsClient.Collection(entity.KindNameImage).Doc(keyName)
	document, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			image = entity.Image{
				Status:    entity.StatusReady,
				CreatedAt: time.Now(),
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
	if _, err := docRef.Set(ctx, &image); err != nil {
		return err
	}
	return nil
}
