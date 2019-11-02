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

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/sugyan/image-dataset/web/entity"
)

type gcp struct {
	csClient   *storage.Client
	dsClient   *datastore.Client
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
	dsClient, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &gcp{
		csClient:   csClient,
		dsClient:   dsClient,
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
	postedAt, err := time.Parse("2006-01-02 15:04:05", data.Meta.PostedAt)
	if err != nil {
		return err
	}
	parts := make([]int, 136)
	for i := 0; i < 68; i++ {
		parts[i*2], parts[i*2+1] = data.Parts[i][0], data.Parts[i][1]
	}
	metaData := struct {
		Angle   float32 `json:"angle"`
		FaceID  int     `json:"face_id"`
		PhotoID int     `json:"photo_id"`
		LabelID int     `json:"label_id"`
	}{
		Angle: data.Angle,
	}
	metaData.FaceID, _ = strconv.Atoi(data.Meta.FaceID)
	metaData.PhotoID, _ = strconv.Atoi(data.Meta.PhotoID)
	metaData.LabelID, _ = strconv.Atoi(data.Meta.LabelID)
	meta, err := json.Marshal(&metaData)
	if err != nil {
		return err
	}

	key := datastore.NameKey(entity.KindNameImage, keyName, nil)
	image := entity.Image{
		ImageURL:  fmt.Sprintf("https://storage.googleapis.com/%s/images/%s", g.bucketName, keyName),
		SourceURL: data.Meta.SourceURL,
		PhotoURL:  data.Meta.PhotoURL,
		Size:      data.Size,
		Size0256:  data.Size >= 256,
		Size0512:  data.Size >= 512,
		Size1024:  data.Size >= 1024,
		Parts:     parts,
		LabelName: data.Meta.LabelName,
		Status:    entity.StatusReady,
		PostedAt:  postedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Meta:      meta,
	}
	if _, err := g.dsClient.Put(ctx, key, &image); err != nil {
		return err
	}
	return nil
}
