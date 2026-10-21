# IceCream POS frontend

## Run locally

1. Start the Go API:

   ```powershell
   cd ..\backend
   go run .
   ```

2. Start the frontend:

   ```powershell
   cd frontend
   npm install
   npm run dev
   ```

The Vite development server proxies `/api` requests to `http://localhost:4000`.
For a separately hosted API, set `VITE_API_URL` to the API URL including `/api`.

## Connected features

- JWT login at `/api/auth/login`
- Dashboard inventory summary
- Inventory list, search, create, update, and delete
- Customer list, search, create, update, and delete
- Role-aware actions for Owner, Manager, and Cashier

The backend seed command creates the Owner, Manager, and Cashier accounts. Set
`SEED_OWNER_PASSWORD`, `SEED_MANAGER_PASSWORD`, and `SEED_CASHIER_PASSWORD`
before running `go run .\cmd\seed` from the `backend` directory. See the
repository [README](../README.md) for the complete setup and API documentation.
