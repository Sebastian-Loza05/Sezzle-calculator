# Decisions

| Decision | Why I took it |
| --- | --- |
| Build the backend in Go. | Go is the company's primary backend language, and the assignment is a good chance to demonstrate it. |
| Use one stateless `POST /calculate` endpoint. | Every operation has the same request and response shape, and calculations do not need storage. |
| Use `net/http` and `encoding/json`. | The standard library handles this small API without extra dependencies. |
| Keep `main`, `httpapi`, and `calculator` separate. | Server startup, HTTP handling, and arithmetic have different responsibilities and can be tested separately. |
| Dispatch through `Calculate` to individual operation functions. | Each formula stays readable while the HTTP handler calls one calculator entry point. |
| Send binary operands as `a` and `b`. | Named fields are simple to read for the calculator's binary operations. |
| Make `sqrt` use only `a` and reject `b`. | Square root is unary, and rejecting an extra value prevents silently ignored input. |
| Define `percentage(a, b)` as `(a / 100) * b`. | “Percentage” is ambiguous without a rule; this means `a` percent of `b`. |
| Allow fractional exponents for positive bases; reject a negative base with a noninteger exponent. | The latter is outside this calculator's real-number result model. |
| Reject negative square roots and `0` raised to a negative power; return `1` for `0^0`. | These rules make the accepted numeric behavior explicit and follow Go's `math.Pow` value for `0^0`. |
| Use `float64` and reject non-finite operands or results. | It suits a general calculator and keeps every successful result valid JSON; exact decimal arithmetic is outside this scope. |
| Validate JSON strictly and return stable error codes. | Missing fields, typos, extra data, and mathematical errors should give predictable feedback to the frontend. |
| Use table-driven unit tests and `httptest`. | They cover arithmetic and HTTP behavior without starting a network server. |
