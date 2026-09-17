# Running Locally

## Architecture note

Only MongoDB runs in Docker (`docker-compose.yml` defines a single `mongodb` service). The Go API itself runs **natively on the host**, not in a container — `db/db.go` hardcodes `mongodb://localhost:27017`, which only resolves correctly when the app is a host process talking to Mongo's published port, not when the app itself is containerized (unless you added it to the compose network and swapped the hostname — not currently set up). A `Dockerfile` exists for building the app as an image, but nothing in `docker-compose.yml` currently runs it.

## Prerequisites

- **Docker** (Desktop or Engine) — runs MongoDB
- **Go 1.23+** — matches `go.mod`; runs the API itself
- **Git** — to clone
- MongoDB Compass (optional, for inspecting data)

## 1. Get the code

```bash
git clone https://github.com/yogisyo16/root-aura-service.git
cd root-aura-service
```

`.env` is committed to the repo with working local defaults — nothing to configure out of the box:
```
BINARY=mongoTodosYogis
DB_CONTAINER_NAME=mongodb-todos-yogis
MONGO_DB=todos_db
MONGO_DB_USERNAME=admin
MONGO_DB_PASSWORD=password
PORT=65510
```

## 2. Start MongoDB

```bash
docker compose up -d
```
(or `docker-compose up -d` / `make up` on a machine with GNU Make — the Makefile target just runs `docker-compose up --build -d --remove-orphans`.)

This starts `mongo:8.0.10` in a container named `mongodb-todos-yogis`, publishing `27017`, seeded with database `todos_db` and root user `admin`/`password` (from `.env`, via Mongo's `MONGO_INITDB_*` variables).

**Compass connection string:**
```
mongodb://admin:password@localhost:27017/todos_db?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false
```

## 3. Run the API

The app reads `PORT`, `MONGO_DB_USERNAME`, `MONGO_DB_PASSWORD` from the environment — there's no `.env`-loading library in the code (`os.Getenv` only), so export them in the shell you run it from.

**macOS / Linux:**
```bash
export PORT=65510 MONGO_DB_USERNAME=admin MONGO_DB_PASSWORD=password
go run ./api
```
or via the Makefile: `make build && make start` (`make restart` does both).

**Windows (PowerShell)** — no `make` by default, so run the equivalent directly:
```powershell
$env:PORT="65510"
$env:MONGO_DB_USERNAME="admin"
$env:MONGO_DB_PASSWORD="password"
go run ./api
```
or build first:
```powershell
go build -o mongoTodosYogis.exe ./api
.\mongoTodosYogis.exe
```

You should see:
```
Successfully connected and pinged MongoDB!
Server is running on port :65510
```

## 4. Verify

```bash
curl http://localhost:65510/api/v1/healthcheck
```
Should return `{"Msg":"Health Check","Code":200}`. See `API.md` for the full endpoint list.

## Stopping

```bash
docker compose down
```
