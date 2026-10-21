# Ice Cream Shop Backend

REST API for ice cream inventory and customer management, built with Go 1.21+
or newer, Fiber v2, GORM, SQLite, JWT authentication, and role-based access
control. A React frontend is included under [frontend](./frontend).

## Requirements

- Go 1.21+
- Node.js 18+ (only needed for the optional React frontend)

## Run the backend

From PowerShell:

```powershell
cd backend
Copy-Item .env.example .env
# Edit .env and replace every placeholder value.
go run .
```

The API runs at `http://localhost:4000`. Health check: `GET /health`.

## Seed users

Set three passwords of at least 12 characters, then run the seed command:

```powershell
$env:SEED_OWNER_PASSWORD = "replace-with-owner-password"
$env:SEED_MANAGER_PASSWORD = "replace-with-manager-password"
$env:SEED_CASHIER_PASSWORD = "replace-with-cashier-password"
go run .\cmd\seed
```

The accounts are `owner@icecream.local`, `manager@icecream.local`, and
`cashier@icecream.local`. Passwords are never stored in the repository.

## API

All inventory and customer endpoints require `Authorization: Bearer <token>`.
Authenticate with `POST /api/auth/login`. Logout is
`POST /api/auth/logout`; because JWT is stateless, clients should discard the
token after logout.

Inventory endpoints use `/api/inventory` and support `?search=<name-or-sku>`.
Customer endpoints use `/api/customers` and support `?search=<name-email-phone>`.
Each inventory item and customer record stores the authenticated creator in
`created_by`, referencing `users.id`.

| Role | Inventory | Customers |
| --- | --- | --- |
| Owner | View, create, edit, delete | View, create, edit, delete |
| Manager | View, create, edit | View, create, edit |
| Cashier | View, edit | View, create, edit |

Import [the Postman collection](./docs/IceCreamShop.postman_collection.json) to
exercise every endpoint. The editable [ER diagram](./docs/er-diagram.mmd)
documents the SQLite schema.

## Run tests

```powershell
cd backend
go test ./...
go vet ./...
```

## Run the frontend

```powershell
cd frontend
npm install
npm run dev
```

The Vite development server proxies `/api` to `http://localhost:4000`.
