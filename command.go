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

type commands struct {
  values map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
  c.values[name] = f;
}

func (c *commands) run(s *state, cmd command) error {
  f, ok := c.values[cmd.name];
  if !ok {
    return errors.New("Could not find command")
  }
  err := f(s, cmd)
  return err
}

func handlerLogin(s *state, cmd command) error {
  l := len(cmd.arguments)
  if l == 0 {
    return errors.New("No login arguments provided (expected 1)")
  } else if l > 1 {
    return errors.New("Too many login arguments provided (expected 1)")
  }

  ctx := context.Background()
  _, err := s.db.GetUser(ctx, cmd.arguments[0])
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
  l := len(cmd.arguments)
  if l == 0 {
    return errors.New("No register arguments provided (expected 1)")
  } else if l > 1 {
    return errors.New("Too many register arguments provided (expected 1)")
  }

  ctx := context.Background()

  _, err := s.db.GetUser(ctx, cmd.arguments[0])
  if err == nil {
    return errors.New("User already in database!")
  }

  user, err := s.db.CreateUser(ctx, database.CreateUserParams{uuid.New(), time.Now(), time.Now(), cmd.arguments[0]})
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
