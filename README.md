# OTP Service (Golang)

This service provides OTP-based login/registration and basic user management with JWT authentication.  
It can run in **two modes**:
- **Postgres + Redis** (persistent, production-like)
- **In-memory** (ephemeral, good for dev/testing)

🗄️ Database Choice
- Postgres → reliable and durable storage for permanent user data (persistence, search, pagination).
- Redis → fast, in-memory store, ideal for temporary data such as OTPs and rate limiting.
✅ Using both ensures a balance of durability and performance.
---

The mode is controlled via `.env` toggles (`USE_DB`, `USE_REDIS`).

- **Rate limiting**
  - Max **3 OTP requests per phone number within 10 minutes**
  - In-memory fixed window implementation

- **JWT Authentication**
  - HS256-signed tokens
  - Expiry: **24h** (configurable)
  - Protects `/users` endpoints

## ⚙️ Features
- OTP-based login & registration
  - OTP stored in Redis or in-memory
  - Rate-limited (3 requests per 10 min per phone)
  - Expires after 2 minutes
- User management
  - List users (with pagination & search)
  - Get user by ID
- JWT-based authentication
- PostgreSQL for user storage (with migration)
- Redis for OTP + rate limiting
- Fallback to in-memory if disabled/unavailable
- Swagger/OpenAPI docs
- Dockerized (multi-stage build with caching)

---

## 📂 Project Structure
```
cmd/server         # main entrypoint
internal/config    # env loading + toggles
internal/domain    # domain entities (User)
internal/infra     # infra (postgres, memory)
internal/otp       # OTP service interfaces + impls
internal/http      # Fiber routing, handlers, middleware
docs/              # generated Swagger docs
```

---

## 🔑 Configuration (.env)

All configuration is centralized in **one `.env` file**. Example:

```env
# ---- API ----
PORT=8080
JWT_SECRET=

# ---- Toggles ----
USE_DB=true
USE_REDIS=true

# ---- Postgres ----
POSTGRES_PORT=5432
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_DNS=db
DATABASE_URL=

# ---- Redis ----
REDIS_PORT=6379
REDIS_DB=
REDIS_DNS=redis
REDIS_URL=

# ---- Tunables ----
# Durations use Go format (e.g., 30s, 2m, 1h, 24h)
OTP_TTL=2m
RATE_LIMIT_MAX=3
RATE_LIMIT_WINDOW=10m
TOKEN_TTL=24h
```
- If `.env` is missing → warning is logged, defaults are used.
- `PORT`: application running port.
- `JWT_SECRET`: the `jwt` secret (sould be set in production).
- If `USE_DB=false` → in-memory user repository. If `true` you should fill in `postgres` data!
- If `USE_REDIS=false` → in-memory OTP/rate limiter. If `true` you should fill in `redis` data!
- If `DATABASE_URL`/`REDIS_URL` are empty but toggles true → URLs are auto-built from base vars.  
- `OTP_TTL`: how long an OTP is valid.
- `RATE_LIMIT_MAX`: how many OTP requests a phone number can make per window.
- `RATE_LIMIT_WINDOW`: sliding window for rate limiting.
- `TOKEN_TTL`: how long JWT tokens remain valid.

---

## 🚀 Running

You can run the service locally with Go, or inside Docker Compose.

### Local (Go)
```bash
make run
```
Runs the server directly with `go run ./cmd/server`

**Note: if you dont run the databases localy and set in `.env` file, then program uses in-memory databses by default!

---

### Docker

#### Build image
```bash
make docker
```
Or just run `docker build -t otp-service:dev .`

#### Build without cache
```bash
make compose-build
```
Or just run `docker compose build --no-cache`


#### Start full stack (detached)
```bash
make up
```
Or just run `docker compose up -d --build`

#### Start full stack (attached logs)
```bash
make up-attached
```
Or just run `docker compose up --build`

#### Start API only (don’t attach db/redis logs)
```bash
make up-api
```
Or just run `docker compose up --build --no-attach db --no-attach redis`

#### Restart quickly (no rebuild)
```bash
make up-no-build
```
Or just run `docker compose up -d`

#### Stop containers (keep database data)
```bash
make down
```
Or just run `docker compose down --remove-orphans`

#### Stop containers and remove all volumes (fresh DB/Redis)
```bash
make down-clean
```
Or just run `docker compose down -v --remove-orphans`

---

## 📝 API Docs

After running, visit:
```
http://localhost:8080/swagger/index.html
```

---

## 🧪 Example Usage

### Request OTP
```bash
curl -X POST http://localhost:8080/api/v1/auth/request-otp   -H 'Content-Type: application/json'   -d '{"phone":"+1555"}'
```

Check logs for OTP code.

### Verify OTP
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp   -H 'Content-Type: application/json'   -d '{"phone":"+1555","otp":"123456"}'
```

Response includes JWT token.

### Get Users
```bash
curl -H "Authorization: Bearer <TOKEN>"   http://localhost:8080/api/v1/users?page=1&size=10
```

---

## 🧩 Development
- Generate Swagger locally:
  ```bash
  make swag
  ```
  Or just run `go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/server/main.go -o ./docs`

- Tidy modules:
  ```bash
  make tidy
  ```
  Or just run `go mod tidy`
---

## 🗄️ Data Inspection

### Postgres
```bash
1. docker exec -it otp_db psql -U $POSTGRES_USER -d $POSTGRES_DB
# inside psql
2. \dt # (optional) shows databases
3. SELECT * FROM users; # shows users
```

### Redis
```bash
1. docker exec -it otp_redis redis-cli
2. keys * # shows all keys
#for each key u can
3. get $key #shows key data
```

---

## 🛡️ Notes
- In-memory mode is ephemeral — users, OTPs, and rate limits vanish on restart.
- For production, always run with Postgres + Redis and set a **strong JWT_SECRET**.
