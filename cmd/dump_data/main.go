package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/sugyan/image-dataset/web/entity"
	"golang.org/x/image/draw"
	"google.golang.org/api/iterator"
)

func main() {
	projectID := flag.String("projectID", "", "project ID")
	size := flag.Int("size", 512, "target image size")
	num := flag.Int("num", 100, "number of dump images")
	flag.Parse()
	if *projectID == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*projectID, *size, *num); err != nil {
		log.Fatal(err)
	}
}

func run(projectID string, size, num int) error {
	urlCh, err := query(context.Background(), projectID, num)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	wg := sync.WaitGroup{}
	for _, w := range newWorkers(10, size) {
		wg.Add(1)
		go func(w *worker) {
			defer wg.Done()
			w.run(urlCh, errCh)
		}(w)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		return err
	}
	return nil
}

func query(ctx context.Context, projectID string, num int) (<-chan string, error) {
	urlCh := make(chan string)

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	go func() {
		defer client.Close()
		query := client.Collection(entity.KindNameImage).
			Where("Size0512", "==", true).
			OrderBy("ID", firestore.Asc).
			Limit(num)
		iter := query.Documents(ctx)
		for i := 0; i < num; i++ {
			document, err := iter.Next()
			if err != nil {
				if errors.Is(err, iterator.Done) {
					break
				} else {
					log.Fatal(err)
				}
			}
			var image entity.Image
			if err := document.DataTo(&image); err != nil {
				log.Fatal(err)
			}
			urlCh <- image.ImageURL
		}
		close(urlCh)
	}()
	return urlCh, nil
}

type worker struct {
	index int
	size  int
}

func newWorkers(num, size int) []*worker {
	workers := []*worker{}
	for i := 0; i < num; i++ {
		workers = append(workers, &worker{
			index: i,
			size:  size,
		})
	}
	return workers
}

func (w *worker) run(urlCh <-chan string, errCh chan<- error) {
	kernel := draw.CatmullRom
	for url := range urlCh {
		log.Printf("[%02d] %s", w.index, url)
		resp, err := http.Get(url)
		if err != nil {
			errCh <- fmt.Errorf("[%s] failed to download: %s", url, err.Error())
			continue
		}
		defer resp.Body.Close()
		// Download and resize
		img, err := jpeg.Decode(resp.Body)
		if err != nil {
			errCh <- fmt.Errorf("[%s] failed to decode: %s", url, err.Error())
			continue
		}
		dst := image.NewRGBA(image.Rect(0, 0, w.size, w.size))
		kernel.Scale(dst, image.Rect(0, 0, w.size, w.size), img, img.Bounds(), draw.Over, nil)
		// Save to file
		file, err := os.Create(filepath.Join("images", fmt.Sprintf("%s.jpg", path.Base(url))))
		if err != nil {
			errCh <- err
			return
		}
		if err := jpeg.Encode(file, dst, &jpeg.Options{Quality: 100}); err != nil {
			errCh <- err
			return
		}
	}
}
