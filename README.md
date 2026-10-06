# Go + Inertia.js v3 Starter Kit

A clean, reusable, production-oriented full-stack starter kit using **Go**, **Inertia.js v3**, **React**, **TypeScript**, **shadcn/ui**, **Tailwind CSS**, and **PostgreSQL**.

## Technology Stack

| Layer     | Technology                          |
|-----------|-------------------------------------|
| Backend   | Go 1.27+                            |
| Frontend  | React 18 + TypeScript (strict mode) |
| SPA Bridge| Inertia.js v3 (gonertia/v3)         |
| UI        | shadcn/ui + Tailwind CSS v3         |
| Build     | Vite 5                              |
| ORM       | GORM v2                             |
| Database  | PostgreSQL 17                       |
| Sessions  | alexedwards/scs                     |
| Auth      | bcrypt password hashing             |

---

## Requirements

- Go 1.21+
- Node.js 18+
- npm 9+
- Docker + Docker Compose (for PostgreSQL)

---

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/thbappy7706/go-inertia-starter-kit.git
cd go-inertia-starter-kit
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Install frontend dependencies

```bash
npm install
```

### 4. Start PostgreSQL with Docker Compose

```bash
docker compose up -d
```

This will start PostgreSQL 17 on port `5432` with:
- **DB**: `starter_kit`
- **User**: `postgres`
- **Password**: `postgres`

### 5. Configure environment

```bash
cp .env.example .env
```

Edit `.env` if you need to customize database credentials or the session secret:

```env
APP_NAME=Go Inertia Starter Kit
APP_ENV=local
APP_PORT=8090
APP_URL=http://localhost:8090

DATABASE_DRIVER=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=starter_kit
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_SSLMODE=disable

SESSION_SECRET=<change-this-to-a-long-random-string>
```

> **Important**: Change `SESSION_SECRET` to a long random string before deploying to production.

---

## Development

### Start the Go backend

```bash
go run ./cmd/server
```

The server starts on: **http://localhost:8090**

The first startup automatically:
1. Runs database migrations (creates `users` and `products` tables)
2. Seeds default admin user and 6 sample products

### Start the Vite frontend dev server (hot reload)

Open a second terminal:

```bash
npm run dev
```

Vite runs on port **5173** and automatically injects HMR into the Go-served pages.

---

## Default Development Login

> **These credentials are for local development ONLY. Never use in production.**

| Field    | Value                  |
|----------|------------------------|
| Email    | admin@example.com      |
| Password | password               |

---

## Database Migrations

Migrations are run automatically on server startup via GORM `AutoMigrate`.

Migration files (reference only) are in:
```
database/migrations/
├── 001_create_users.sql
└── 002_create_products.sql
```

---

## Seed Data

The seeder runs automatically on startup (if no data exists). To re-seed, connect to the database and truncate the tables, then restart the server.

Seeded data includes:
- 1 admin user (`admin@example.com` / `password`)
- 6 sample products with descriptions, pricing, and status variations

---

## Project Structure

```
go-inertia-starter-kit/
├── cmd/
│   └── server/
│       └── main.go              # Application entrypoint
│
├── internal/
│   ├── config/                  # Environment/config loading
│   ├── database/                # GORM connection + AutoMigrate
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Auth, Guest, ShareProps middleware
│   ├── models/                  # GORM models + response DTOs
│   ├── repositories/            # Data access layer
│   ├── services/                # Business logic
│   ├── session/                 # Session manager + FlashProvider
│   └── testhelper/              # Shared test utilities
│
├── routes/
│   └── routes.go                # Chi router + route definitions
│
├── database/
│   ├── migrations/              # Reference SQL migration files
│   └── seeders/                 # Development seed data
│
├── resources/
│   └── views/
│       └── app.html             # Root Inertia HTML template
│
├── public/
│   ├── hot                      # Vite HMR hot file (dev only)
│   └── build/                   # Vite production build output
│       ├── manifest.json
│       └── assets/
│
├── frontend/
│   ├── main.tsx                 # React + Inertia entrypoint
│   ├── css/app.css              # Tailwind + CSS variables
│   ├── components/
│   │   ├── ui/                  # shadcn/ui component library
│   │   ├── DataTable.tsx
│   │   ├── Pagination.tsx
│   │   ├── EmptyState.tsx
│   │   ├── LoadingState.tsx
│   │   ├── ConfirmDialog.tsx
│   │   ├── PageHeader.tsx
│   │   └── StatusBadge.tsx
│   ├── layouts/
│   │   ├── AppLayout.tsx        # Authenticated admin layout
│   │   ├── AuthLayout.tsx       # Login/Register layout
│   │   └── components/
│   │       ├── Sidebar.tsx
│   │       ├── Header.tsx
│   │       └── UserMenu.tsx
│   ├── lib/
│   │   ├── utils.ts             # cn(), formatCurrency(), formatDate()
│   │   └── validations/
│   │       ├── auth.ts          # loginSchema, registerSchema (Zod)
│   │       └── product.ts       # productSchema (Zod)
│   ├── pages/
│   │   ├── Auth/
│   │   │   ├── Login.tsx
│   │   │   └── Register.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Products/
│   │   │   ├── Index.tsx
│   │   │   ├── Create.tsx
│   │   │   ├── Edit.tsx
│   │   │   └── components/
│   │   │       └── ProductForm.tsx
│   │   └── Users/
│   │       └── Index.tsx
│   └── types/
│       └── index.ts             # Shared TypeScript types
│
├── docker-compose.yml           # PostgreSQL 17 service
├── .env                         # Local environment (not in git)
├── .env.example                 # Environment template
├── vite.config.ts               # Vite configuration
├── tailwind.config.js           # Tailwind CSS configuration
├── tsconfig.json                # TypeScript configuration
└── package.json                 # Frontend dependencies
```

---

## Features

### Authentication
- Registration with full validation (name, email, password, confirmation)
- Login with "remember me"
- Secure bcrypt password hashing (DefaultCost)
- HTTP-only session cookies (SCS session manager)
- Session renewal on login (prevents session fixation)
- Inertia shared props: `auth.user` available on every page
- `RequireAuth` middleware protects all dashboard routes
- `RequireGuest` middleware redirects authenticated users from login/register

### Dashboard
- Welcome message with user name
- Real-time stats: total users, total products, active products
- Recent products table with status badges and edit links

### Products (Full CRUD)
- Server-side paginated listing (10 per page)
- Search by name, slug, description
- Filter by status (active / inactive)
- URL-synchronized filters (`?search=...&status=...&page=...`)
- Create with auto-generated slug
- Edit with validation
- Delete with AlertDialog confirmation
- Toast notifications for all operations
- Server-side validation with Zod client-side validation

### Users
- Paginated user listing
- Search by name or email

### UI
- shadcn/ui components throughout
- Responsive sidebar (Sheet/drawer on mobile)
- Sonner toast notifications
- Loading states, empty states, error states
- Accessible semantic HTML

---

## Available Scripts

| Script                   | Description                           |
|--------------------------|---------------------------------------|
| `go run ./cmd/server`    | Start Go backend                      |
| `npm run dev`            | Start Vite HMR dev server             |
| `npm run build`          | Production frontend build             |
| `npm run typecheck`      | TypeScript type checking              |
| `npm run lint`           | ESLint linting                        |
| `npm run format`         | Prettier formatting                   |
| `go test ./...`          | Run all backend tests                 |
| `go fmt ./...`           | Format Go code                        |
| `go vet ./...`           | Go static analysis                    |

---

## Production Build

1. Build the frontend:
   ```bash
   npm run build
   ```
2. Build the Go binary:
   ```bash
   go build -o ./server ./cmd/server
   ```
3. Set production environment:
   ```env
   APP_ENV=production
   SESSION_SECRET=<very-long-random-secret>
   DATABASE_SSLMODE=require
   ```
4. Run the server:
   ```bash
   ./server
   ```

---

## Tests

Backend tests cover:
- ✅ Registration success
- ✅ Registration validation failures (name, email, password)
- ✅ Duplicate email rejection
- ✅ Password hashing verification
- ✅ Login success
- ✅ Invalid credentials rejection
- ✅ Product CRUD (create, update, delete, find)
- ✅ Product validation failures
- ✅ Product search, filter by status, pagination
- ✅ User listing, search, pagination
- ✅ Protected route requires authentication
- ✅ Guest route redirects authenticated users
- ✅ Logout destroys session

Run tests:
```bash
go test -v ./internal/services/... ./routes/...
```

---

## Architecture Notes

- **No JWT** — Uses secure HTTP-only session cookies via `alexedwards/scs`
- **No REST API** — Frontend communicates via Inertia.js (server-rendered page props)
- **FlashProvider** — Custom implementation bridging SCS sessions to Gonertia flash data
- **Soft deletes** — GORM `gorm.DeletedAt` on User and Product models
- **ShareAuth middleware** — Runs on every request, populates `auth.user` + `flash` Inertia shared props
- **Vite integration** — Uses `gonertia/v3` built-in Vite support with HMR hot file detection

---

## Security

- Passwords hashed with bcrypt (`DefaultCost`)
- HTTP-only session cookies
- Session renewal on login (prevents session fixation attacks)
- `Secure` cookie flag enabled in production (`APP_ENV=production`)
- SQL injection protection via GORM parameterized queries
- No password field exposed to React/Inertia via `UserResponse` DTO
- No secrets exposed to frontend
- Proper HTTP status codes throughout