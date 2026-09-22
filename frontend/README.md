# Calculator frontend

React and TypeScript frontend built with Vite and Tailwind CSS.

## Run locally

Start the Go backend on its default port, `8080`, then run:

```sh
npm install
npm run dev
```

Open the URL printed by Vite. Its development proxy sends `/calculate` requests
to `http://localhost:8080`. If the backend uses another port, change the proxy
target in `vite.config.ts`.

## Verify

```sh
npm test
npm run lint
npm run build
```

The UI validates required finite operands. It always displays the second
operand; for square root, the field is disabled and its value is omitted from
the request. Backend messages are shown for calculation errors.
