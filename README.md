# gator
RSS feed aggregator from the guided project on boot.dev

A CLI tool that uses:
- a [postgreSQL](https://www.postgresql.org/) database to store user information
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) to generate Go code from SQL queries so that my code can interact with the database.

Currently handles commands:
- ```register <name>```: adds a user with <name> to the database
- ```login <name>```: sets current user to <name>
- ```reset```: removes all users from the database (DANGER: THIS IS IRREVERSIBLE)
- ```users```: prints a list of all users in the database
