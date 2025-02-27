package main

import (
  "internal/config"
  "internal/database"
)

type state struct {
  db *database.Queries
  cfg *config.Config
}
