# gator
RSS feed aggregator from the guided project on boot.dev

A CLI tool that uses:
- a [postgreSQL](https://www.postgresql.org/) database to store user information
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) to generate Go code from SQL queries so that my code can interact with the database.

Currently handles commands:
- ```register <name>```: adds a user with ```<name>``` to the database
- ```login <name>```: sets current user to ```<name>```
- ```reset```: removes all users from the database (DANGER: THIS IS IRREVERSIBLE)
- ```users```: prints a list of all users in the database
- ```agg```: currently fetches the feed for ```https://www.wagslane.dev/index.xml```
- ```addfeed <name> <url>```: adds the feed with ```<name>``` found at ```<url>``` to the current user's feeds
- ```feeds```: prints a list of all feeds in the database
- ```follow <url>```: follows the feed at ```<url>``` for the current user
- ```following```: prints a list of all feeds that the current user is following
- ```unfollow <url>```: unfollows the feed at ```<url>``` for the current user

Uses 3 tables:
- ```users```: contains information on each user
- ```feeds```: contains information on each feed
- ```feed_follows```: a joining table containing information on which feeds a user is following

## NOTES

### Migrations

Write up and down migrations in the sql/schema folder. This will be a ```.sql``` file containing the sql code that Goose will use for each migration. Each file should have both an up and corresponding down migration that reverses the changes (if possible). Goose will use these files to make database changes.

For the up migration, navigate to sql/schema and run:
```goose postgres <connection_string> up```
(```goose postgres postgres://carolinemillan:@localhost:5432/gator up```)

For the down migration, navigate to sql/schema and run:
```goose postgres <connection_string> down```
(```goose postgres postgres://carolinemillan:@localhost:5432/gator down```)

NOTE: connection string is "postgres://carolinemillan:@localhost:5432/gator"

### Queries

Write queries in the sql/queries folder. This will be a ```.sql``` file containing the sql code that SQLC will use to generate Go code to interact with the database. Make a new file for each table in the database.

To generate the Go code from the queries, run:
```sqlc generate```
