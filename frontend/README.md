# JobTrack frontend

This folder contains the React, TypeScript, and Vite application shell.

## Requirements

- Node.js 20.19+ or 22.12+
- npm

## Run locally

From this folder:

```sh
npm install
cp .env.example .env.local
npm run dev
```

Vite prints the local development URL in the terminal (usually
`http://localhost:5173`). The app currently provides navigable placeholder
pages for Dashboard, Applications, and Keyword Matcher.

## API configuration

Set `VITE_API_BASE_URL` in `.env.local` to the backend origin, for example
`http://localhost:8080`. The value is available to frontend code through
`src/config.ts`. Vite variables are bundled into browser code, so do not put
secrets in them.

## Other commands

```sh
npm run build
npm run preview
```
