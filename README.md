# Sezzle Calculator

A full-stack calculator project with a Go backend and React, TypeScript, Vite,
and Tailwind CSS frontend.

## Run the full stack with Docker

From the repository root:

```sh
docker compose up --build
```

Open `http://localhost:3000`. The frontend container serves the built React app
with Nginx and forwards `/calculate` to the Go container. The backend is only
reachable on the Compose network. Stop both services with `docker compose down`.

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

Send a JSON POST request to `/calculate`. These examples use the backend run
directly on port `8080`; with Docker Compose, use `localhost:3000` instead:

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
{"error":{"code":"division_by_zero","message":"cannot divide by zero"}}
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
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go vet ./...
```

From `frontend/`:

```sh
npm test
npm run coverage
npm run lint
npm run build
```

See [COVERAGE.md](COVERAGE.md) for the measured coverage of both layers. Go's
`coverage.out` and Vitest's `frontend/coverage/` HTML report are generated
locally and ignored by Git.


## Decisions

| Decision | Why I took it |
| --- | --- |
| Build the backend in Go. | I chose Go because it is the company's primary backend language, and I wanted this project to show that I can build and test a small service in it. |
| Use one stateless `POST /calculate` endpoint. | I used one endpoint because every operation accepts the same kind of request and returns one result. There is no calculation history to store, so separate endpoints or a database would add work without helping the user. |
| Use `net/http` and `encoding/json`. | I stayed with Go's standard library because it already provides the routing and JSON handling this API needs. That keeps the service easy to build and avoids a framework dependency for one route. |
| Keep `main`, `httpapi`, and `calculator` separate. | I put startup, HTTP rules, and arithmetic in separate packages so I can test calculations without HTTP and test request handling without starting a server. |
| Dispatch through `Calculate` to individual operation functions. | I gave the HTTP handler one function to call, then kept each arithmetic rule in a small named function. This makes the operations easier for me to read and test. |
| Send binary operands as `a` and `b`. | I chose named fields because a request such as `{"a":10,"b":5}` makes each operand's role clear, especially for subtraction, division, and percentage. |
| Make `sqrt` use only `a` and reject `b`. | I treated square root as a one-operand operation. Rejecting `b` catches a mistaken request instead of accepting and silently ignoring part of it. |
| Define `percentage(a, b)` as `(a / 100) * b`. | I wrote down the meaning because “percentage” could mean several calculations. Here `a` is the percent and `b` is the amount, so `20` and `50` produce `10`. |
| Allow fractional exponents for positive bases; reject a negative base with a noninteger exponent. | I wanted useful results such as `9^0.5 = 3`. For negative bases, a decimal `float64` exponent does not reliably express which fractional powers have real results, so I reject those requests consistently. |
| Reject negative square roots and `0` raised to a negative power; return `1` for `0^0`. | I made these edge cases explicit so users receive a clear error where there is no finite real result. I kept Go's `math.Pow` behavior for `0^0` instead of inventing a separate rule. |
| Use `float64` and reject non-finite operands or results. | I chose ordinary floating-point numbers for a general calculator, including fractional inputs and exponents. I reject infinities and `NaN` because they cannot be returned as valid JSON numbers; exact money arithmetic is outside this task. |
| Validate JSON strictly and return stable error codes. | I reject missing fields, unknown fields, and extra JSON so typos do not look like valid calculations. Error codes also give the frontend a predictable response while messages explain the problem to users. |
| Use table-driven unit tests and `httptest`. | I used tables to check many inputs against the same arithmetic rules, and `httptest` to check status codes, valid results, and malformed JSON without running a separate server. |
| Build the frontend with React, TypeScript, Vite, and Tailwind CSS. | I kept the existing Vite setup, used TypeScript for typed requests and operation names, and chose Tailwind because I am comfortable using it and even the AI is building it I need to understand what it's doing. |
| Keep the single-screen UI in `App.tsx` with local state; separate validation and HTTP calls into `calculator.ts` and `api.ts`. | I kept form state next to the only component that uses it. I moved input rules and `fetch` into small files so I can test them directly and keep `App.tsx` focused on interaction and display. |
| Keep the second operand visible but disabled for square root, and omit `b` from its request. | I kept the form layout stable when the operation changes, then disabled the field to show it is unused. Omitting `b` also matches the backend's square-root contract. |
| Validate required finite operands in the UI and show backend error messages. | I check simple input mistakes before sending a request so users get quick feedback. I still show server errors because only the backend decides whether a mathematical operation is valid. |
| Proxy `/calculate` to Go during Vite development. | I let the browser use the same relative URL in development and deployment. The Vite proxy reaches Go locally, so I do not need to change the backend just to handle browser CORS. |
| Abort calculations after 15 seconds. | I added a limit so a stalled request does not leave the form stuck in its loading state. The user gets a timeout message and can try again. |
| Test validation and API behavior with Vitest. | I tested zero, missing and non-finite inputs, the square-root request shape, server errors, malformed responses, network failures, and timeouts because these are the paths most likely to break a small calculator. |
| Test the rendered UI with Vitest and jsdom. | I checked the actual form behavior: square root keeps the second field visible but disabled, and validation or backend failures appear on screen. These tests catch UI regressions that function tests alone would miss. |
| Build separate Go and frontend images with Compose; serve the frontend and proxy `/calculate` through Nginx. | I kept each image responsible for one application layer. Nginx serves the built React files and forwards API calls to Go on the Compose network, so the browser uses one origin. |
| Report Go and frontend test coverage in `COVERAGE.md`. | I included the measured results because the take-home asks for coverage, and I documented the commands so a reviewer can reproduce them. I leave generated coverage files out of Git to keep the repository small. |
| Prioritize useful tests over reaching 100% coverage. | I added tests for malformed requests and responses because those fallbacks matter to users. `main.go` and `main.tsx` mainly wire up startup, so I checked the running backend with a Docker smoke test and the frontend bundle with a production build.