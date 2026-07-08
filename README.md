## Gator
Gator is a cli tool that lets you add and follow RSS feeds. Multiple users can be registered. The users and their different follows are saved in a local Postgres database

This project was made for the [Boot.dev](http://boot.dev/) back-end development path. It was a guided course where I dealt with most implementation details and learned to use sql in a project for the first time, using tools such as goose and sqlc

## Installation
You will need both Postgres, Go and Goose to install, setup and run Gator

For goose, once you have go installed just run
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Once you have them installed, you'll need to do a git clone to get the code
```
git clone https://github.com/Robot-tim1/gator.git
```
Go into the newly created directory
```
cd gator
```
When you have it run go install so you can run the program anywhere
```
go install
```
Now onto the database part. Once you've installed postgres and set your password if needed depending on your os enter the psql shell. On linux it's
```
sudo -u postgres psql
```
Then create and setup the gator database
```
CREATE DATABASE gator;
\c gator
ALTER USER postgres PASSWORD 'postgres';
```
You can make the password whatever you want, but I used postgres for simplicity

Now to make sure the database is setup we'll have to run a migration with Goose. You'll need your connection string, which can be different depending on os, but on linux it looks something like this
```
postgres://postgres:postgres@localhost:5432/gator
```
The format being 'protocol://username:password@host:port/database' so if you set it up differently you'll need it to be different. Make sure you're in the projects root directory and run this with your own connection string
```
goose -dir ./sql/schema postgres postgres://postgres:postgres@localhost:5432/gator up
```
Then you'll have to make a config file in your home directory, typically ~/.gatorconfig.json. It should include something like this if you have the same connection string, otherwise use your own plus '?sslmode=disable'
```
{
  "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```
Now that should be all the steps to set it up. If you want to see the commands you can run, look below to see available commands  

## Commands
Prefix all commands with gator

| Command | Description |
|---------|-------------|
| `register <username>` | Registers and logs in as given name |
| `login <username>` | Logs into any registered users |
| `users` | Displays all registered users and currently logged in user |
| `addfeed <name> <url>` | Adds feed with a custom name |
| `feeds` | Displays all added feeds and which user added them |
| `follow <url>` | Follows a feed by the given url |
| `unfollow <url>` | Unfollows a feed by the given url |
| `following` | Displays all the feeds the user follows by name |
| `agg [time]` | Aggregates from followed feeds in intervals of the given time (e.g., 1m, 30s). default 10 seconds |
| `browse [limit]` | Displays the posts aggregated from the feeds, post number limited by given number, default 2 |
| `reset` | Deletes all the data from the database, resetting everything. Use at your own discretion |

Use agg in the background on a different terminal, while you interact with the cli on another terminal
