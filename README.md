# Sezzle Calculator

A full-stack calculator project with a Go backend and React, TypeScript, Vite,
and Tailwind CSS frontend.

## Run the backend

The backend currently uses the Go version declared in [backend/go.mod](backend/go.mod).
From the repository root:

```sh
cd backend
go run ./cmd/server
```

The server listens on port `8080` by default. Set `PORT` to use another port:

```sh
PORT=8081 go run ./cmd/server
```

## Run the frontend

In a second terminal, from the repository root:

```sh
cd frontend
npm install
npm run dev
```

Open the local URL printed by Vite. The development server proxies `/calculate`
to the backend on port `8080`. Run `npm test`, `npm run lint`, and
`npm run build` from `frontend/` to verify the frontend. If you use a different
backend port, update the proxy target in `frontend/vite.config.ts`.

## API

Send a JSON POST request to `/calculate`:

```sh
curl -i -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":10,"b":2}'
```

The response is `{"result":5}`. Supported operations are `add`, `subtract`,
`multiply`, `divide`, `power`, `sqrt`, and `percentage`. The binary operations
require both `a` and `b`, including when either is zero. `sqrt` takes only `a`:

```sh
curl -i -X POST http://localhost:8080/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":9}'
```

For `percentage`, `a` means the percentage and `b` is the base amount. For
example, `{"operation":"percentage","a":20,"b":50}` returns `{"result":10}`.

Errors have a stable JSON shape. For example, division by zero returns HTTP 400:

```json
{"error":{"code":"division_by_zero","message":"division by zero"}}
```

The backend also rejects malformed JSON, missing fields, unsupported operations,
unknown fields, oversized bodies, non-finite results, negative square roots,
negative bases with fractional exponents, and zero raised to a negative power. See
[backend/ARCHITECTURE.md](backend/ARCHITECTURE.md) for the API rules and design
rationale. [DECISIONS.md](DECISIONS.md) records the key implementation choices.

## Tests and coverage

From `backend/`:

```sh
go test ./...
go test -cover ./...
go vet ./...
```

Current `go test -cover ./...` results:

| Package | Statement coverage |
| --- | ---: |
| `internal/calculator` | 100.0% |
| `internal/httpapi` | 90.0% |
| `cmd/server` | 0.0% (startup wiring) |

The tests cover the calculation rules and the HTTP request/response contract.
The startup command is intentionally small; run it with the API example above
to check it as a complete service.
