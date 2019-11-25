package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	projectID := flag.String("projectID", "", "project ID")
	datadir := flag.String("datadir", "", "data directory")
	flag.Parse()
	if *projectID == "" || *datadir == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(*projectID, *datadir); err != nil {
		log.Fatal(err)
	}
	log.Println("finish")
}

func run(projectID, datadir string) error {
	pathsCh, err := walk(datadir)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	wg := sync.WaitGroup{}
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			worker(index, projectID, pathsCh, errCh)
		}(i)
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		log.Println(err)
	}

	return nil
}

func walk(datadir string) (<-chan string, error) {
	pathsCh := make(chan string)
	go func() {
		defer close(pathsCh)
		if err := filepath.Walk(datadir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}
			pathsCh <- path
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}()
	return pathsCh, nil
}

func worker(index int, projectID string, pathsCh <-chan string, errCh chan<- error) {
	gcp, err := newGcp(projectID)
	if err != nil {
		errCh <- err
		return
	}
	for filepath := range pathsCh {
		log.Printf("[%02d] %s", index, filepath)
		if err := gcp.upload(filepath); err != nil {
			errCh <- fmt.Errorf("error [%s]: %s", filepath, err.Error())
		}
	}
}
