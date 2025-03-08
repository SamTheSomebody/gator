package main

import (
  "errors"
  "fmt"
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

  err := argumentCountError(len(cmd.arguments), saved.targetArgumentCount)
  if err != nil {
    return err
  }

  err = saved.f(s, cmd)
  return err
}

func argumentCountError(count int, target int) error {
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
