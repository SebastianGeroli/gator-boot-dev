# Gator

Gator is a command-line RSS feed aggregator. It scrapes feeds you follow on a
schedule and stores the posts in a Postgres database, which you can then
browse from the CLI.

## Requirements

Before running Gator you'll need:

- **Postgres** (v15+ recommended) — Gator stores users, feeds, and posts in
  a Postgres database.
- **Go** (v1.21+) — used to build/install the CLI.

## Installing

Install the `gator` CLI with `go install`:

```bash
go install github.com/SebastianGeroli/gator-boot-dev/cmd/gator@latest
```

This downloads, compiles, and installs the `gator` binary to your
`$GOPATH/bin` (usually `~/go/bin`) — make sure that directory is on your
`PATH`. Since Go produces statically compiled binaries, once it's installed
you can run `gator` directly without needing the Go toolchain around.

## Configuration

Gator reads its configuration from a `.gatorconfig.json` file in your home
directory. Create one at `~/.gatorconfig.json` with the following shape:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "user_name": ""
}
```

- `db_url` — connection string for your Postgres database. Update it to
  match your local setup.
- `user_name` — leave this blank; it's managed automatically once you
  `register` or `login`.

Make sure the database and schema exist (run the migrations in
`sql/schema` with [goose](https://github.com/pressly/goose)) before using
the CLI.

## Running

Once installed and configured, run commands with:

```bash
gator <command> [args...]
```

A few commands to get started with:

- `gator register <name>` — create a new user and log in as them.
- `gator login <name>` — switch the current user.
- `gator addfeed <name> <url>` — add a new RSS feed and automatically follow it.
- `gator follow <url>` — follow an existing feed as the current user.
- `gator following` — list the feeds the current user follows.
- `gator agg <interval>` — continuously scrape feeds for new posts (e.g. `gator agg 1m`).
- `gator browse [limit]` — print the most recent posts (defaults to 2).
