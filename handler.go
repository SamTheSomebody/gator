package main

import (
  "errors"
  "fmt"
  "github.com/google/uuid"
  "time"
  "context"
  "internal/database"
  "strconv"
)

func handlerLogin(s *state, cmd command) error {
  _, err := s.db.GetUser(context.Background(), cmd.arguments[0])
  if err != nil {
    return err
  }

  err = s.cfg.SetUser(cmd.arguments[0])
  if err != nil {
    return err
  }

  fmt.Println("Username has been set")
  return nil
}

func handlerRegister(s *state, cmd command) error {
  _, err := s.db.GetUser(context.Background(), cmd.arguments[0])
  if err == nil {
    return errors.New("User already in database!")
  }

  params := database.CreateUserParams{uuid.New(), time.Now(), time.Now(), cmd.arguments[0]}
  user, err := s.db.CreateUser(context.Background(), params)
  if err != nil {
    return err
  }

  err = s.cfg.SetUser(user.Name)
  if err != nil {
    return err 
  }
  fmt.Printf("Created new user! Name: %v, Created At: %v, Updated At: %v, ID: %v", user.Name, user.CreatedAt, user.UpdatedAt, user.ID)
  return nil
}

func handlerReset(s *state, cmd command) error {
  err := s.db.DeleteUsers(context.Background())
  if err != nil {
    return err
  }
  
  fmt.Println("Successfully reset user table")
  return nil
}

func handlerListUsers(s *state, cmd command) error {
  users, err := s.db.GetUsers(context.Background())
  if err != nil {
    return err
  }

  for _, x := range(users) {
    output := "* " + x.Name
    if x.Name == s.cfg.CurrentUserName {
      output += " (current)"
    }
    fmt.Println(output)
  }

  return nil
}

func handlerAggregate(s *state, cmd command) error {
  timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
  if err != nil {
    return err
  }

  fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
  
  ticker := time.NewTicker(timeBetweenRequests)
  for ; ; <-ticker.C {
    scrapeFeeds(s)
  }
  return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
  params := database.CreateFeedParams{uuid.New(), time.Now(), time.Now(), cmd.arguments[0], cmd.arguments[1], user.ID}
  feed, err := s.db.CreateFeed(context.Background(), params)
  if err != nil {
    return err
  }
  fmt.Println(feed)

  cmd.arguments = cmd.arguments[1:]
  f := middlewareLoggedIn(handlerFollowFeed)
  return f(s, cmd)
}

func handlerGetFeeds(s *state, cmd command) error {
  feeds, err := s.db.GetFeeds(context.Background())
  if err != nil {
    return err
  }

  for _, feed := range feeds {
    user, err := s.db.GetUserByID(context.Background(), feed.UserID)
    if err != nil {
      return err
    }
    fmt.Printf("%v: %v (%v)\n", feed.Name, feed.Url, user.Name)
  }

  return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
  feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
  if err != nil {
    return err
  }

  params := database.CreateFeedFollowParams{uuid.New(), time.Now(), time.Now(), user.ID, feed.ID}
  feed_follow, err := s.db.CreateFeedFollow(context.Background(), params)
  if err != nil {
    return err
  }

  fmt.Println(feed_follow)
  return nil
}

func handlerListFollowingFeeds(s *state, cmd command, user database.User) error {
  feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
  if err != nil {
    return err
  }

  for _, feed := range feeds {
    fmt.Println(feed)
  }
  return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
  feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
  if err != nil {
    return err
  }

  params := database.RemoveFeedFollowParams{user.ID, feed.ID}
  return s.db.RemoveFeedFollow(context.Background(), params) 
}

func handlerBrowse(s *state, cmd command, user database.User) error {
  limit := 2
  if len(cmd.arguments) == 1 {
    i, err := strconv.Atoi(cmd.arguments[0])
    if err != nil {
      return err
    }
    limit = i
  }

  params := database.GetPostsForUserParams{user.ID, int32(limit)}
  posts, err := s.db.GetPostsForUser(context.Background(), params)
  if err != nil {
    return err
  }
  
  fmt.Println("Retrieving posts...\n")

  for _, x := range posts {
    fmt.Printf("Post from: %v\n", x.FeedName)
    fmt.Printf("Title: %v\n", x.Title)
    fmt.Printf("Link: %v\n", x.Url)
    fmt.Printf("Published: %v\n", x.PublishedAt)
    fmt.Printf("Description: %v\n", x.Description)
    fmt.Println("")
  }

  fmt.Printf("Retreived %v posts (requested %v)\n", len(posts), limit)
  return nil
}
