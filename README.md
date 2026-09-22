# City counter

Type a letter, get the number of cities (out of the 50 most populous, per
[GeoNames](http://api.geonames.org/searchJSON?featureClass=P&maxRows=50&orderby=population&username=hsample))
whose name starts with it. E.g. `c` → 3 (Cairo, Chengdu, Chongqing), `r` → 1 (Rio de Janeiro).

- **Backend** – Go 1.27, standard library only (`net/http`). Fetches the GeoNames list, caches it in memory, exposes one endpoint.
- **Frontend** – React 19 + TypeScript (Vite), Tailwind v4, [shadcn/ui](https://ui.shadcn.com) components, TanStack Query for data fetching.

## Run with Docker Compose (recommended)

```sh
docker compose up --build
```

Open <http://localhost:3000>. nginx serves the built frontend and proxies `/api/*` to the Go backend.

Optional env vars (set in the shell or a `.env` file next to `docker-compose.yml`):

| Variable            | Default   | Purpose                           |
| ------------------- | --------- | --------------------------------- |
| `GEONAMES_USERNAME` | `hsample` | GeoNames account used for lookups |
| `CACHE_TTL`         | `10m`     | How long the city list is cached  |

## Run locally without Docker

Backend (needs Go ≥ 1.22):

```sh
cd backend
go run .            # listens on :8080
```

Frontend (needs Node ≥ 20):

```sh
cd frontend
npm install
npm run dev         # http://localhost:5173, proxies /api to :8080
```

## Tests

```sh
cd backend && go test ./...
```

Covers the letter matching (case-insensitive, Unicode-aware), the GeoNames client (parsing, caching within TTL, refetch after TTL, errors not cached) and the HTTP handler (happy path, validation, upstream failure, CORS).

## API

`GET /api/cities/count?letter=<single letter>`

```json
{ "letter": "c", "count": 3, "cities": ["Chengdu", "Cairo", "Chongqing"] }
```

- `400` – `letter` missing, longer than one character, or not a letter
- `502` – GeoNames unreachable or returned an error


