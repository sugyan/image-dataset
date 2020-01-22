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

var (
	projectID string
	size      int
	num       int
	status    string
	outdir    string
)

func init() {
	flag.StringVar(&projectID, "projectID", "", "project ID")
	flag.IntVar(&size, "size", 512, "target image size")
	flag.IntVar(&num, "num", 100, "number of dump images")
	flag.StringVar(&status, "status", "", "target status")
	flag.StringVar(&outdir, "outdir", "images", "path to output directory")
}

func main() {
	flag.Parse()
	if projectID == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	urlCh, err := query(context.Background())
	if err != nil {
		return err
	}
	errCh := make(chan error)
	wg := sync.WaitGroup{}
	for _, w := range newWorkers(20) {
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

func query(ctx context.Context) (<-chan string, error) {
	urlCh := make(chan string)

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	go func() {
		defer client.Close()
		query := client.Collection(entity.KindNameImage).
			Where("Size0512", "==", true).
			OrderBy("ID", firestore.Asc)
		if status != "" {
			switch status {
			case entity.StatusReady.Path():
				query = query.Where("Status", "==", entity.StatusReady)
			case entity.StatusNG.Path():
				query = query.Where("Status", "==", entity.StatusNG)
			case entity.StatusPending.Path():
				query = query.Where("Status", "==", entity.StatusPending)
			case entity.StatusOK.Path():
				query = query.Where("Status", "==", entity.StatusOK)
			default:
				log.Fatalf("invalid status: %s", status)
			}
		}
		i := 0
	Loop:
		for {
			log.Printf("%d", i)
			iter := query.Limit(500).Documents(ctx)
			for {
				document, err := iter.Next()
				if err != nil {
					if errors.Is(err, iterator.Done) {
						break
					} else {
						log.Fatal(err)
					}
				}
				query = query.StartAfter(document)

				var image entity.Image
				if err := document.DataTo(&image); err != nil {
					log.Fatal(err)
				}
				urlCh <- image.ImageURL
				i++
				if i == num {
					break Loop
				}
			}

		}
		close(urlCh)
	}()
	return urlCh, nil
}

type worker struct {
	index int
}

func newWorkers(numWorkers int) []*worker {
	workers := []*worker{}
	for i := 0; i < numWorkers; i++ {
		workers = append(workers, &worker{
			index: i,
		})
	}
	return workers
}

func (w *worker) run(urlCh <-chan string, errCh chan<- error) {
	outdir, err := filepath.Abs(outdir)
	if err != nil {
		errCh <- err
		return
	}
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
		dst := image.NewRGBA(image.Rect(0, 0, size, size))
		kernel.Scale(dst, image.Rect(0, 0, size, size), img, img.Bounds(), draw.Over, nil)
		// Save to file
		file, err := os.Create(filepath.Join(outdir, fmt.Sprintf("%s.jpg", path.Base(url))))
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
