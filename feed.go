package main

import (
  "net/http"
  "html"
  "io"
  "encoding/xml"
  "context"
)

type RSSFeed struct {
  Channel struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    Item []RSSItem `xml:"item"`
  } `xml:"channel"`
}

type RSSItem struct {
  Title string `xml:"title"`
  Link string `xml:"link"`
  Description string `xml:"description"`
  PubDate string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
  req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
  if err != nil {
    return nil, err
  }
  req.Header.Set("User-Agent", "gator")

  client := http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  data, err := io.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  feed := &RSSFeed{}
  err = xml.Unmarshal(data, feed)
  if err != nil {
    return nil, err
  }

  feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
  feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
  for i, x := range feed.Channel.Item {
    feed.Channel.Item[i].Title = html.UnescapeString(x.Title)
    feed.Channel.Item[i].Description = html.UnescapeString(x.Description)
  }

  return feed, nil
}
