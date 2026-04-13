# gator
RSS feed aggregator from the guided project on boot.dev

A CLI tool that uses:
- a [postgreSQL](https://www.postgresql.org/) database to store user information
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) to generate Go code from SQL queries so that my code can interact with the database.
