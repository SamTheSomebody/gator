package main

import (
  "database/sql"
  "log"
  "os"

  "internal/database"
	"internal/config"

  _ "github.com/lib/pq"
)

type state struct {
  db *database.Queries
  cfg *config.Config
}

func main() {
  cfg, err := config.Read()
  if err != nil {
    log.Fatal(err)
  }

  db, err := sql.Open("postgres", cfg.DBURL)
  if err != nil {
    log.Fatal(err)
  }

  dbQueries := database.New(db)

  s := state {
    db: dbQueries,
    cfg: &cfg,
  }

  commands := commands {
    values: make(map[string]savedCommand),
  }

  commands.register("login", handlerLogin, 1)
  commands.register("register", handlerRegister, 1)
  commands.register("reset", handlerReset, 0)
  commands.register("users", handlerListUsers, 0)
  commands.register("agg", handlerAggregate, 1)
  commands.register("addfeed", middlewareLoggedIn(handlerAddFeed), 2)
  commands.register("feeds", handlerGetFeeds, 0)
  commands.register("follow", middlewareLoggedIn(handlerFollowFeed), 1)
  commands.register("following", middlewareLoggedIn(handlerListFollowingFeeds), 0)
  commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed), 1)
  c := createCommand()

  err = commands.run(&s, c)
  if err != nil {
    log.Fatal(err)
  }
}

func createCommand() command {
  if len(os.Args) < 2 {
    log.Fatal("[Fatal] Less than 2 arguments provided")
  }

  name := os.Args[1]

  var args []string
  if len(os.Args) > 2 {
    args = os.Args[2:]
  }

  c := command{
    name: name,
    arguments: args,
  }

  return c
}
