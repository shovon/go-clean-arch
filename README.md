# Simple Hexagonal/Clean architecture example

A demo of how to structure Go code to adhere to the hexagonal architecture, making it easier to test code.

Also demostrates a dependency injection technique.

## Environment variables that MUST be set

- `PGCONNSTRING`: a URL to a PostgreSQL database

## Getting started

Things that you will need installed:

- Go
- Docker
- PostgreSQL

Be sure to have PostgreSQL running, and be sure you know what URL to use to connect to it.

Then, invoke this command (remember to set the environment variables!):

```bash
./tasks/migrate.sh up
```

Then run the app with `go run .`. But remember to set the environment variables!
