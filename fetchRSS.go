package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
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
				// fetch posts from feed
				log.Printf("Fetching posts from feed at %s...", url)
				channel, err := fetchRSS(feed.Url)
				if err != nil {
					log.Printf("error fetching feed at url: %s\nerror: %s", url, err)
					return
				} else {
					err = apiConfig.markFeedFetched(feed)
					if err != nil {
						log.Printf("error updating last_fetched_at for feed at %s:  %s", feed.Url, err)
						return
					}
				}
				for _, item := range channel.Items {
					// check that item pubDate can be converted to time.Time
					pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
					if err != nil {
						log.Printf("could not parse post pubDate: %s", item.PubDate)
						pubDate = time.Now().UTC()
					}

					// only get posts since the last time posts were fetched from the feed
					var lastFetchedAt time.Time
					if feed.LastFetchedAt == nil {
						aYear := 24 * 365 * time.Hour
						lastFetchedAt = time.Now().Add(-1 * aYear).UTC()
					} else {
						lastFetchedAt = *feed.LastFetchedAt
					}
					if pubDate.After(lastFetchedAt) {
						// create post parameters and insert to database
						postParams := database.CreatePostParams{
							ID:          uuid.New(),
							CreatedAt:   time.Now().UTC(),
							UpdatedAt:   time.Now().UTC(),
							Title:       item.Title,
							Url:         item.Link,
							Description: item.Description,
							PublishedAt: pubDate,
							FeedID:      feed.ID,
						}
						_, err = apiConfig.DB.CreatePost(context.TODO(), postParams)
						if err != nil {
							log.Printf("error creating post from feed at %s: %s", feed.Url, err.Error())
						}
					}
				}
			}(feed.Url)
		}
		waitGroup.Wait()
	}
}
