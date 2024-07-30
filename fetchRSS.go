package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/nicpatlan/RSS_Feed_Aggregator/internal/database"
)

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func fetchRSS(url string) (Channel, error) {
	// fetch the url
	res, err := http.Get(url)
	if err != nil {
		return Channel{}, err
	}
	defer res.Body.Close()

	// decode the body into a go struct
	channel := new(Channel)
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(channel)
	if err != nil {
		return Channel{}, err
	}
	return *channel, nil

	/*
		// TODO remove : testing responses pubDates
		//layout := "Mon, 02 Jan 2006 15:04:05 Z0700" same as time.RFC1123Z
		t, err := time.Parse(time.RFC1123Z, channel.Items[0].PubDate)
		if err != nil {
			log.Fatalf("time.Parse error: %s", err)
		}
		log.Printf("the item is %v", channel.Items[0])
		log.Printf("the parsed time is: %s", t)
	*/
}

func (apiConfig *apiConfig) getNextFeedsToFetch(n int32) ([]Feed, error) {
	// get least recently fetched feeds from database
	dbFeeds, err := apiConfig.DB.GetNextFeedsToFetch(context.TODO(), n)
	if err != nil {
		return make([]Feed, 0), err
	}
	return convertDatabaseFeedsToArray(dbFeeds), nil
}

func (apiConfig *apiConfig) markFeedFetched(feed Feed) error {
	markFeedParams := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: feed.ID,
	}
	err := apiConfig.DB.MarkFeedFetched(context.TODO(), markFeedParams)
	if err != nil {
		log.Printf("error updating last_feed_fetched: %s", err)
		return err
	}
	return nil
}

func (apiConfig *apiConfig) fetchFeedBatch(batchSize int32, fetchInterval time.Duration) {
	// set and start fetch timing
	ticker := time.NewTicker(fetchInterval)
	for ; ; <-ticker.C {
		var waitGroup sync.WaitGroup
		feeds, err := apiConfig.getNextFeedsToFetch(batchSize)
		if err != nil {
			log.Printf("error fetching feeds from database: %s", err)
		}
		for _, feed := range feeds {
			// increment goroutine wait group counter
			waitGroup.Add(1)
			go func(url string) {
				// decrement goroutine wait group counter when done
				defer waitGroup.Done()
				channel, err := fetchRSS(feed.Url)
				if err != nil {
					log.Printf("error fetching feed at url: %s\nerror: %s", feed.Url, err)
					return
				} else {
					err = apiConfig.markFeedFetched(feed)
					if err != nil {
						log.Printf("error updating last_fetched_at for feed at %s:  %s", feed.Url, err)
						return
					}
				}
				for _, item := range channel.Items {
					log.Printf("%s feed with title: %s", feed.Name, item.Title)
				}
			}(feed.Url)
		}
		waitGroup.Wait()
	}
}
