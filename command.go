package main

import (
  "errors"
  "fmt"
  "github.com/google/uuid"
  "time"
  "context"
  "internal/database"
)

type command struct {
  name string
  arguments []string
}

type savedCommand struct {
  f func(*state, command) error
  targetArgumentCount int
}
type commands struct {
  values map[string]savedCommand
}

func (c *commands) register(name string, f func(*state, command) error, targetArgumentCount int) {
  c.values[name] = savedCommand{f, targetArgumentCount}
}

func (c *commands) run(s *state, cmd command) error {
  saved, ok := c.values[cmd.name];
  if !ok {
    return errors.New("Could not find command")
  }

  err := generateArgumentCountError(len(cmd.arguments), saved.targetArgumentCount)
  if err != nil {
    return err
  }

  err = saved.f(s, cmd)
  return err
}

func generateArgumentCountError(count int, target int) error {
  if count == target {
    return nil
  }
  text := "Too "
  if count > target {
    text += "many"
  } else {
    text += "few"
  }
  text += fmt.Sprintf(" arguments provided! (received %v, expected %v)", count, target)
  return errors.New(text)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
  return func(s *state, cmd command) error {
    user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
      return err
    }
    return handler(s, cmd, user)
  }
}

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
  err := s.db.Reset(context.Background())
  if err != nil {
    return err
  }
  
  fmt.Println("Successfully reset user table")
  return nil
}

func handlerGetUsers(s *state, cmd command) error {
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
  url := "https://www.wagslane.dev/index.xml"
  feed, err := fetchFeed(context.Background(), url)
  if err != nil {
    return err
  }
 
  fmt.Println(feed)
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

func handlerGetFollowingFeeds(s *state, cmd command, user database.User) error {
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
  params := database.RemoveFeedFollowParams{user.ID, cmd.arguments[0]}
  return s.db.RemoveFeedFollow(context.Background(), params) 
}
