package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/sugyan/image-dataset/web/entity"
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

func (g *gcp) delete(ctx context.Context, images []*entity.Image) error {
	batch := g.fsClient.Batch()
	for _, image := range images {
		// delete from storage
		log.Printf("image %s: (size: %d)", image.ID, image.Size)
		obj := g.csClient.Bucket(g.bucketName).Object(fmt.Sprintf("images/%s", image.ID))
		if err := obj.Delete(ctx); err != nil {
			log.Printf("failed to delete object: %v", obj)
			// continue to delete
		}
		// delete from firestore
		docRef := g.fsClient.Collection(entity.KindNameImage).Doc(image.ID)
		batch.Delete(docRef)
	}
	if _, err := batch.Commit(ctx); err != nil {
		return err
	}
	return nil
}
