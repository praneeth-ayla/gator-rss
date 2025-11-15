# Gator CLI

Gator is a small command line tool that works with user accounts, feeds, and feed follows. It uses Postgres, goose for migrations, and sqlc for database code.

## Requirements

- Go installed
- PostgreSQL installed
- goose installed
- sqlc installed

## Install

Clone the repo:

```
git clone https://github.com/praneeth-ayla/gator-rss
cd gator-rss
```

Install the binary:

```
go install
```

Now you can run:

```
gator
```

## Database Setup

Create the database:

```
createdb gator
```

Run migrations with goose:

```
goose -dir sql/schema postgres "postgres://localhost:5432/gator?sslmode=disable" up
```

Generate sqlc code if needed:

```
sqlc generate
```

## Config File

Create this file:

```
~/.gatorconfig.json
```

Example:

```json
{
  "db_url": "postgres://localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Running the Program

Development:

```
go run .
```

Normal usage:

```
gator
```

## Basic Commands

Register a user:

```
gator register alice
```

Log in:

```
gator login alice
```

Add a feed:

```
gator addfeed https://example.com/feed.xml
```

Follow a feed:

```
gator follow https://example.com/feed.xml
```

See follows:

```
gator following
```

Scrape feeds:

```
gator agg <time>
```

## Project Layout

```
internal/config
internal/database
sql/queries
sql/schema
```

- goose runs the `sql/schema` files
- sqlc reads `sql/queries` and generates Go code
- config manages your CLI config
- handlers and commands are in the root folder
