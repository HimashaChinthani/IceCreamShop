# Ice Cream Shop

Backend API and optional React frontend for managing an ice cream shop's
inventory and customer records.

## Technology

- Go 1.21 or newer
- Fiber v2
- GORM
- SQLite
- JWT authentication
- React and Vite frontend

The backend is the required part of the assessment. The React application is
included as an optional client for demonstrating the API. No WPF application is
included.

## Project structure

```text
.
├── backend/
│   ├── cmd/seed/              # Initial user seed command
│   ├── database/              # SQLite connection and migrations
│   ├── handlers/              # Authentication, inventory, and customer APIs
│   ├── middleware/            # JWT authentication and role authorization
│   ├── models/                # GORM models
│   └── routes/                # Protected API routes
├── docs/
│   ├── er-diagram.mmd         # Editable Mermaid ER diagram
│   └── IceCreamShop.postman_collection.json
└── frontend/                  # Optional React/Vite client
```

## Prerequisites

Install the following:

- Go 1.21 or newer
- Node.js 18 or newer (only required for the frontend)
- npm (included with Node.js)

Verify the installations:

```powershell
go version
node --version
npm --version
```

## Backend setup

All commands in this section are run from the repository root unless stated
otherwise.

### 1. Configure environment variables

Copy the environment template:

```powershell
Copy-Item .\backend\.env.example .\backend\.env
```

Open `backend\.env` and replace the placeholder values. `JWT_SECRET` must be
at least 32 bytes long. Each seed password must be at least 12 characters.

Example:

```env
JWT_SECRET=use-a-long-random-secret-value-at-least-32-bytes
SEED_OWNER_PASSWORD=owner-password-at-least-12
SEED_MANAGER_PASSWORD=manager-password-at-least-12
SEED_CASHIER_PASSWORD=cashier-password-at-least-12
```

Do not commit `backend\.env`, real passwords, or the SQLite database. They are
ignored by Git.

### 2. Download dependencies

```powershell
Set-Location .\backend
go mod download
```

### 3. Create the database and seed users

The application automatically creates and migrates `backend\icecream.db` when
it starts. Seed the three required roles before logging in:

```powershell
go run .\cmd\seed
```

The seed command creates these accounts:

| Role | Email |
|---|---|
| Owner | `owner@icecream.local` |
| Manager | `manager@icecream.local` |
| Cashier | `cashier@icecream.local` |

Passwords are read from the `SEED_*_PASSWORD` variables in `backend\.env`.
Existing users are not overwritten when the seed command is run again.

### 4. Start the API

From `backend\`:

```powershell
go run .
```

The API listens on:

```text
http://localhost:4000
```

Health check:

```text
GET http://localhost:4000/health
```

Expected response:

```json
{
  "status": "healthy"
}
```

Keep this terminal running while using the frontend or Postman.

## Authentication

Login with:

```http
POST http://localhost:4000/api/auth/login
Content-Type: application/json
```

```json
{
  "email": "owner@icecream.local",
  "password": "<the-owner-password-from-backend-.env>"
}
```

The response contains the JWT in `access_token`:

```json
{
  "access_token": "<jwt>",
  "token_type": "Bearer",
  "expires_in": 43200,
  "user": {
    "role": "OWNER"
  }
}
```

Use the token on every protected request:

```http
Authorization: Bearer <jwt>
```

The token expires after 12 hours. Login again after it expires.

Logout:

```http
POST http://localhost:4000/api/auth/logout
Authorization: Bearer <jwt>
```

JWT logout is stateless. The server validates the request, and the client must
discard the token. The frontend performs this local token removal when the
Logout button is used.

## API endpoints

All inventory and customer endpoints require a valid Bearer token.

### Authentication

| Method | Endpoint | Access |
|---|---|---|
| POST | `/api/auth/login` | Public |
| POST | `/api/auth/logout` | Authenticated users |
| GET | `/api/auth/me` | Authenticated users |

### Inventory

Inventory records contain a name, category/flavor, SKU, quantity, and price.
Search by name or SKU with the optional `search` query parameter.

| Method | Endpoint | Owner | Manager | Cashier |
|---|---|---:|---:|---:|
| GET | `/api/inventory` | Yes | Yes | Yes |
| GET | `/api/inventory/:id` | Yes | Yes | Yes |
| POST | `/api/inventory` | Yes | Yes | No |
| PUT | `/api/inventory/:id` | Yes | Yes | Yes |
| DELETE | `/api/inventory/:id` | Yes | No | No |

Examples:

```text
GET  /api/inventory
GET  /api/inventory?search=chocolate
GET  /api/inventory/1
POST /api/inventory
PUT  /api/inventory/1
DELETE /api/inventory/1
```

Create or update body:

```json
{
  "name": "Chocolate",
  "category": "Chocolate",
  "sku": "CHO-001",
  "quantity": 20,
  "price": 3.5
}
```

### Customers

Customer records can be searched by name, email, or phone.

| Method | Endpoint | Owner | Manager | Cashier |
|---|---|---:|---:|---:|
| GET | `/api/customers` | Yes | Yes | Yes |
| GET | `/api/customers/:id` | Yes | Yes | Yes |
| POST | `/api/customers` | Yes | Yes | Yes |
| PUT | `/api/customers/:id` | Yes | Yes | Yes |
| DELETE | `/api/customers/:id` | Yes | No | No |

Examples:

```text
GET  /api/customers
GET  /api/customers?search=alice
GET  /api/customers/1
POST /api/customers
PUT  /api/customers/1
DELETE /api/customers/1
```

Create or update body:

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "phone": "0771234567",
  "address": "Colombo"
}
```

Requests without a valid token return `401 Unauthorized`. Authenticated users
without permission for an operation return `403 Forbidden`.

## Postman

Import [the Postman collection](./docs/IceCreamShop.postman_collection.json).

Before running protected requests:

1. Start the backend.
2. Run the Login request.
3. Copy the `access_token` value from the response.
4. Set the collection variable `token` to that JWT.
5. Use `Authorization: Bearer {{token}}` for protected requests.

The collection variables are:

| Variable | Default |
|---|---|
| `baseUrl` | `http://localhost:4000` |
| `token` | Empty; set after login |
| `itemId` | `1` |
| `customerId` | `1` |

Use the seeded account matching the operation you want to demonstrate. For
example, use the cashier account to verify that inventory creation and
deletion return `403 Forbidden`.

## Run the tests

From `backend\`:

```powershell
go test ./...
go vet ./...
```

## Optional React frontend

The frontend consumes the API and provides login, dashboard, inventory, and
customer screens with role-aware actions.

Open a second terminal and run:

```powershell
Set-Location .\frontend
npm install
npm run dev
```

Open the URL printed by Vite, normally:

```text
http://localhost:5173
```

The Vite development server proxies `/api` requests to
`http://localhost:4000`. To use a separately hosted API, create
`frontend\.env` with:

```env
VITE_API_URL=https://your-api-host.example/api
```

Frontend checks:

```powershell
npm run lint
npm run build
```

## Database and ER diagram

SQLite is stored in `backend\icecream.db`. It is created automatically by
GORM migrations when the backend starts. The editable ER diagram is available
at [docs/er-diagram.mmd](./docs/er-diagram.mmd) and describes users,
inventory, and customer relationships.

To recreate the local database during development, stop the backend and
delete only the local database file:

```powershell
Remove-Item .\backend\icecream.db
```

Then start the backend and run the seed command again.

## Troubleshooting

### `JWT_SECRET must contain at least 32 bytes`

Set a longer `JWT_SECRET` in `backend\.env`. Make sure the file is named
exactly `.env` and is inside the `backend` directory.

### Login returns `Invalid email or password`

Run the seed command from `backend\` and verify that the password used for
login matches the corresponding `SEED_*_PASSWORD` value.

### Protected requests return `401`

Check that the request includes:

```http
Authorization: Bearer <jwt>
```

Do not include quotation marks around the token.

### Requests return `403`

The authenticated role does not have permission for that operation. Check the
role matrix above and log in with the appropriate seeded account.

### The frontend cannot connect to the API

Confirm that the backend is running on port `4000`, then refresh the frontend.
If using a custom API host, check the `VITE_API_URL` value and restart Vite
after changing the environment file.
