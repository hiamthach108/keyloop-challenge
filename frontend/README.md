# Keyloop Inventory Frontend

Scenario B inventory dashboard initialized from the Dreon Next.js boilerplate.

## Stack

- Next.js 16
- React 19
- TypeScript 6
- Ant Design 6 with App Router registry
- Tailwind CSS 4
- ESLint flat config and Prettier 3
- Storybook 10
- Husky and lint-staged

## Getting Started

```sh
npm run dev
```

Prefer the root project commands for normal development:

```sh
make setup
make dev
```

## Data source

The dashboard defaults to API mode and communicates with the Go backend.

```env
NEXT_PUBLIC_INVENTORY_SOURCE=api
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Set `NEXT_PUBLIC_INVENTORY_SOURCE=mock` for isolated UI development. In mock
mode, manager actions are persisted in browser local storage.

## Local commands

```sh
npm install
npm run dev
```

## Scripts

```sh
npm run dev
npm run build
npm run start
npm run lint
npm run typecheck
npm run format
npm run storybook
npm run build-storybook
```

## Notes

This template uses npm and commits `package-lock.json` for reproducible installs.
