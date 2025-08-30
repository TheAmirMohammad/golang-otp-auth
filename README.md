# OTP Service (Golang)

This service provides OTP-based login/registration and basic user management with JWT authentication.  
It can run in **two modes**:
- **Postgres + Redis** (persistent, production-like)
- **In-memory** (ephemeral, good for dev/testing)

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
JWT_SECRET=golangotpauthentication

# ---- Toggles ----
USE_DB=true
USE_REDIS=true

# ---- Postgres ----
POSTGRES_USER=otp
POSTGRES_PASSWORD=otp
POSTGRES_DB=otp
POSTGRES_PORT=5432
POSTGRES_DNS=db
DATABASE_URL=

# ---- Redis ----
REDIS_PORT=6379
REDIS_DB=0
REDIS_DNS=redis
REDIS_URL=
```

- If `USE_DB=false` → in-memory user repository.  
- If `USE_REDIS=false` → in-memory OTP/rate limiter.  
- If `DATABASE_URL`/`REDIS_URL` are empty but toggles true → URLs are auto-built from base vars.  
- If `.env` is missing → warning is logged, defaults are used.  

---

## 🚀 Running

You can run the service locally with Go, or inside Docker Compose.

### Local (Go)
```bash
make run
```
Runs the server directly with `go run ./cmd/server`.

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

#### Stop and clean everything
```bash
make down
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
  Or just run `docker compose up --build`

- Tidy modules:
  ```bash
  make tidy
  ```
- Run with Docker:
  ```bash
  make up
  ```

---

## 🗄️ Data Inspection

### Postgres
```bash
docker exec -it otp_db psql -U $POSTGRES_USER -d $POSTGRES_DB
# inside psql
\dt
SELECT * FROM users;
```

### Redis
```bash
docker exec -it otp_redis redis-cli
keys *
#for each row u can
get $row
```

---

## 🛡️ Notes
- In-memory mode is ephemeral — users, OTPs, and rate limits vanish on restart.
- For production, always run with Postgres + Redis and set a **strong JWT_SECRET**.
