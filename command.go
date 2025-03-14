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
  minArgumentCount int
  maxArgumentCount int
}

type commands struct {
  values map[string]savedCommand
}

func (c *commands) register(name string, f func(*state, command) error, minArgumentCount int, maxArgumentCount int) {
  c.values[name] = savedCommand{f, minArgumentCount, maxArgumentCount}
}

func (c *commands) run(s *state, cmd command) error {
  saved, ok := c.values[cmd.name];
  if !ok {
    return errors.New("Could not find command")
  }

  err := argumentCountError(len(cmd.arguments), saved)
  if err != nil {
    return err
  }

  err = saved.f(s, cmd)
  return err
}

func argumentCountError(count int, cmd savedCommand) error {
  if count >= cmd.minArgumentCount && count <= cmd.maxArgumentCount{
    return nil
  }
  text := "Too "
  if count >= cmd.maxArgumentCount{
    text += "many"
  } else {
    text += "few"
  }
  text += fmt.Sprintf(" arguments provided! (received %v, expected ", count)
  if cmd.minArgumentCount == cmd.maxArgumentCount {
    text += fmt.Sprintf("%v)", cmd.minArgumentCount)
  } else {
    text += fmt.Sprintf("%v - %v)", cmd.minArgumentCount, cmd.maxArgumentCount)
  }
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
