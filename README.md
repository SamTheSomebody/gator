##About the project

Gator is a RSS feed aggre*gator*. It was built as a guided project for a course on [boot.dev]. It was built with Go & Postgres with help from [goose](https://github.com/pressly/goose) & [SQLC](https://sqlc.dev/)

##Getting started

Here's all you need to get up and running

###Prerequisites
**ProgresSQL v15 or later**
**Go**

###Installation
After downloading the gator file, run the following command:
```
go install <download_path>/gator
```
You'll also need to create a config file named *.gatorconfig.json* in your root directory. It must contain:
```
{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}
```
You may need to adjust your database url depending on your operating system and postgres login details.

Alternatively you can build from source by downloading this repository and running:
```
go build <download_path>
go install <download_path>/gator
```

##Commands

*Login: set username to be the active user
```
gator login <username>
```
*Register: add a new username to the database and set it as the active user
```
gator register <username>
```
*Reset: delete all entries from the database
```
gator reset
```
*Users: list all registered usernames
```
gator users
```
*Aggregate: scrape the RSS feeds followed by user with a chosen delay (must be valid input for go time's [ParseDuration](https://pkg.go.dev/time#ParseDuration) i.e. *10s*). **AVOID DDOSing** by setting the delay to a reasonable length! You can terminal this command at point (typically *Crtl+C* or *Cmd+C*)
```
gator agg <duration>
```
*Add Feed: add a feed to the database with a name and url, the active user will automatically follow this feed
```
gator addfeed <feed name> <feed url>
```
*Feeds: list all registered feeds
```
gator feeds
```
*Follow: add a registered feed to the active user's follow list
```
gator follow <feed url>
```
*Following: list all of the active user's followed feeds
```
gator following
```
*Unfollow: remove a registered feed from the active user's follow list
```
gator unfollow <feed url>
```
*Browse: retrieve a number of most recently published posts (default is 2)
```
gator browse <optional:limit>
```
