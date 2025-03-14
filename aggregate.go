package main

import (
  "context"
  "log"
  "github.com/google/uuid"
  "time"
  "internal/database"
  "database/sql"
  "fmt"
  "strings"
)

func scrapeFeeds(s *state) {
  feed, err := s.db.GetNextFeedToFetch(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  _, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
  if err != nil {
    log.Fatal(err)
  }

  rss_feed, err := fetchFeed(context.Background(), feed.Url)
  if err != nil {
    log.Fatal(err)
  }

  savePosts(s, rss_feed, feed.ID)
}

func savePosts(s *state, feed *RSSFeed, feedID uuid.UUID) {
  if feed == nil {
    log.Fatal("Feed is empty!\n")
  }

  for _, item := range feed.Channel.Item {
    pubDate, err := time.Parse("Mon, 02 Jan 2006 03:04:05 -0700", item.PubDate)
    if err != nil {
      fmt.Println("Error: ", err)
    }

    nPubDate := sql.NullTime{pubDate, true}
    nDesc := sql.NullString{item.Description, true}
    if len(item.Description) == 0 {
      nDesc.Valid = false
    }
    params := database.CreatePostParams{uuid.New(), time.Now(), time.Now(), item.Title, item.Link, nDesc, nPubDate, feedID}
    _, err = s.db.CreatePost(context.Background(), params)
    urlErr :=  "duplicate key value violates unique constraint \"posts_url_key\""
    if err != nil && strings.Contains(err.Error(), urlErr) == false {
      fmt.Println(err)
    }
  }
  fmt.Printf("Saved posts from %v\n", feed.Channel.Title)
}
