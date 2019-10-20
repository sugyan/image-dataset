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
	"strings"

	"cloud.google.com/go/storage"
)

type gcp struct {
	csClient   *storage.Client
	bucketName string
}

func newGcp(projectID string) (*gcp, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucketName := projectID + ".appspot.com"
	if os.Getenv("DEVELOPMENT") != "" {
		bucketName = "staging." + bucketName
	}
	return &gcp{
		csClient:   client,
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
	return nil
}
