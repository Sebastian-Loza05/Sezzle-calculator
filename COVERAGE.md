# Coverage report

| Layer | Statement coverage | Other coverage | Scope |
| --- | ---: | --- | --- |
| Go backend | 86.0% | — | All Go packages, including server startup. |
| React frontend | 94.36% (67/71) | 94.64% branches (53/56); 92.30% functions (12/13); 96.61% lines (57/59) | All application TS/TSX files in `src/`, excluding tests. |

Backend package coverage: `internal/calculator` 100.0%, `internal/httpapi`
91.4%, and `cmd/server` 0.0%. Frontend statement coverage: `calculator.ts`
100%, `App.tsx` 93.33%, `api.ts` 95.45%, and `main.tsx` 0%. The uncovered
`main` files are startup wiring; the tests focus on calculation, request, and UI
behavior.

Regenerate from the repository root:

```sh
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

cd ../frontend
npm ci
npm run coverage
```

The frontend command writes a detailed HTML report to
`frontend/coverage/index.html` and a machine-readable summary to
`frontend/coverage/coverage-summary.json`.
