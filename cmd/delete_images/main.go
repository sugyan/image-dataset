package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/sugyan/image-dataset/web/entity"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := flag.String("projectID", "", "project ID")
	flag.Parse()
	if *projectID == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(context.Background(), *projectID); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, projectID string) error {
	gcp, err := newGcp(projectID)
	if err != nil {
		return err
	}
	query := gcp.fsClient.Collection(entity.KindNameImage).
		Where("Size0256", "==", false).
		OrderBy("ID", firestore.Asc)
Loop:
	for {
		iter := query.Limit(200).Documents((ctx))
		images := []*entity.Image{}
		for {
			document, err := iter.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				} else {
					return err
				}
			}
			query = query.StartAfter(document)

			var image entity.Image
			if err := document.DataTo(&image); err != nil {
				return err
			}
			images = append(images, &image)
		}
		if len(images) == 0 {
			break Loop
		}
		if err := gcp.delete(ctx, images); err != nil {
			return err
		}
		log.Printf("deleted %d images", len(images))
	}
	return nil
}
