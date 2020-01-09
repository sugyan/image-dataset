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
	if err := run(*projectID); err != nil {
		log.Fatal(err)
	}
	log.Println("finish")
}

func run(projectID string) error {
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer fsClient.Close()

	// count all images
	stats := []*entity.Count{
		&entity.Count{},
		&entity.Count{},
		&entity.Count{},
	}
	query := fsClient.Collection(entity.KindNameImage).Query
	iter := query.Documents(ctx)
	i := 0
	for {
		document, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			} else {
				return err
			}
		}
		var image entity.Image
		if err := document.DataTo(&image); err != nil {
			return err
		}
		for i, b := range []bool{image.Size0256, image.Size0512, image.Size1024} {
			if b {
				switch image.Status {
				case entity.StatusReady:
					stats[i].Ready++
				case entity.StatusNG:
					stats[i].NG++
				case entity.StatusPending:
					stats[i].Pending++
				case entity.StatusOK:
					stats[i].OK++
				}
			}
		}
		i++
		if i%5000 == 0 {
			log.Printf("%d...", i)
		}
	}
	// update stats
	for i, c := range stats {
		docID := []string{"0256", "0512", "1024"}[i]
		ref := fsClient.Collection(entity.KindNameCount).Doc(docID)
		if _, err := ref.Set(ctx, c); err != nil {
			return err
		}
	}
	return nil
}
