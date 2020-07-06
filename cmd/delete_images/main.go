package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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
		Where("Size0512", "==", false).
		OrderBy("ID", firestore.Asc)

	for {
		if err := gcp.fsClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
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
				return fmt.Errorf("Finised")
			}

			counts := []map[entity.Status]int{{}, {}, {}}
			for _, image := range images {
				for i, size := range []int{256, 512, 1024} {
					if image.Size >= size {
						counts[i][image.Status]++
					}
				}
			}
			if err := gcp.delete(ctx, images); err != nil {
				return err
			}
			log.Printf("deleted %d images", len(images))
			for i, count := range counts {
				if len(count) > 0 {
					docID := []string{"0256", "0512", "1024"}[i]
					ref := gcp.fsClient.Collection(entity.KindNameCount).Doc(docID)
					updates := []firestore.Update{}
					for k, v := range count {
						updates = append(updates, firestore.Update{
							Path:  k.Path(),
							Value: firestore.Increment(-v),
						})
					}
					if err := tx.Update(ref, updates); err != nil {
						return err
					}
					log.Printf("updated counts for %s: %v", docID, count)
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
}
