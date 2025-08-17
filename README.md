# rest-api-go-gin

Minimal notes for running and creating database migrations used by this project.

## Prerequisites

- Go 1.24+
- (Optional) Air for live reload during development: https://github.com/cosmtrek/air or https://github.com/air-verse/air
- golang-migrate CLI for creating/applying migrations: https://github.com/golang-migrate/migrate

## Install tools (examples, Windows PowerShell)

- Install golang-migrate (binary):

  ```powershell
  scoop install migrate  # if using scoop
  # or download from https://github.com/golang-migrate/migrate/releases
  ```

- Install Air (optional):
  ```powershell
  scoop install air
  # or follow project README
  ```

## Create new migrations

From project root, create migration SQL files in the migrations folder:

```powershell
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_users_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_events_table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_attendees_table
```

This creates `*.up.sql` and `*.down.sql` files under `cmd/migrate/migrations`.

## Run migrations

From the project root (so relative paths match), run the migrate command implemented in this repo:

- Up (apply all up migrations):

```powershell
go run ./cmd/migrate up
```

- Down (apply all down migrations / rollback):

```powershell
go run ./cmd/migrate down
```

Notes:

- The program opens `./data.db` by default and reads migrations from `cmd/migrate/migrations`. Run commands from the repository root or adjust paths accordingly.
- The SQLite driver must be available at build time (commonly `github.com/mattn/go-sqlite3` or an alternative pure-Go driver `modernc.org/sqlite`). If you prefer a pure-Go driver to avoid cgo, replace the driver in the code and update imports.

## Module / dependency tips

- To clean unused indirect modules:

```powershell
go mod tidy
```

- To check why a module is present:

```powershell
go mod why -m github.com/mattn/go-sqlite3
```
