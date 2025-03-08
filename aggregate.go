package main

import (
  "context"
  "log"
)

func scrapeFeeds(s *state) {
  feed, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  err = s.db.MarkFeedFetched(context.Background(), feed.ID)
  if err != nil {
    log.Fatal(err)
  }

  rss_feed, err := fetchFeed(context.Background(), feed.Url)
  if err != nil {
    log.Fatal(err)
  }
  
  printFeed(rss_feed)
}
