package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	screenName := flag.String("screen_name", "", "target screen_name")
	flag.Parse()
	if *screenName == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := run(*screenName); err != nil {
		log.Fatal(err)
	}
}

func run(screenName string) error {
	var (
		consumerKey    = os.Getenv("CONSUMER_KEY")
		consumerSecret = os.Getenv("CONSUMER_SECRET")
	)
	config := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	ids := map[int64]bool{}
	maxID := int64(0)
	for i := 0; i < 16; i++ {
		params := &twitter.UserTimelineParams{
			ScreenName:      screenName,
			Count:           200,
			MaxID:           maxID,
			TrimUser:        twitter.Bool(true),
			ExcludeReplies:  twitter.Bool(true),
			IncludeRetweets: twitter.Bool(false),
		}
		tweets, resp, err := client.Timelines.UserTimeline(params)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.New(resp.Status)
		}
		for _, tweet := range tweets {
			if tweet.ExtendedEntities == nil {
				continue
			}
			for _, media := range tweet.ExtendedEntities.Media {
				if _, exist := ids[media.ID]; exist {
					continue
				}
				columns := []string{
					media.IDStr,
					media.MediaURLHttps,
					fmt.Sprintf("https://twitter.com/%s/status/%s", screenName, tweet.IDStr),
					tweet.CreatedAt,
					tweet.User.IDStr,
					screenName,
				}
				fmt.Printf("%v\n", strings.Join(columns, "\t"))

				ids[media.ID] = true
			}
			maxID = tweet.ID
		}
		time.Sleep(time.Second)
	}

	return nil
}
