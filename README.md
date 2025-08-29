# OTP Service (Golang)

A clean-architecture backend service that supports **OTP-based login & registration**, **JWT authentication**, and **basic user management**.  
This README covers everything up to **Phase B** (in-memory repository + Docker build).  

---

## ✨ Features
- **OTP-based login & registration**
  - Secure 6-digit OTP, cryptographically random
  - OTP printed to console (no SMS integration yet)
  - Expires after **2 minutes**
  - One-time use (removed after successful validation)

- **Rate limiting**
  - Max **3 OTP requests per phone number within 10 minutes**
  - In-memory fixed window implementation

- **JWT Authentication**
  - HS256-signed tokens
  - Expiry: **24h** (configurable)
  - Protects `/users` endpoints

- **User management**
  - Get user by ID
  - List users (with pagination & search by phone)

- **In-memory user repository**
  - No DB required
  - Data lost on restart (simple demo mode)

- **Swagger/OpenAPI documentation**
  - Live docs at `http://localhost:8080/swagger/index.html`

- **Dockerfile**
  - Multi-stage build (Go → distroless)
  - Runs with **in-memory repo**

---

## 📂 Project Structure
```
otp-service/
├─ cmd/server/            # Entrypoint
│   └─ main.go
├─ internal/
│   ├─ config/            # Config loader (env-based)
│   ├─ domain/user/       # User model + repository interface
│   ├─ infra/memory/      # In-memory repo implementation
│   ├─ http/
│   │   ├─ handlers/      # Fiber handlers (auth, users)
│   │   └─ router.go      # Route registration + JWT middleware
│   ├─ jwt/               # JWT utilities
│   ├─ otp/               # OTP manager
│   └─ rate/              # Rate limiter
├─ docs/                  # Swagger JSON/YAML (generated)
├─ Dockerfile             # Multi-stage container build
├─ Makefile               # Common tasks (run, swag, tidy)
├─ go.mod / go.sum
└─ README.md
```

---

## ⚙️ Requirements
- **Go 1.25+**
- **Make** (optional, for convenience)
- **Docker** 

---

## 🏃 Running Locally (in-memory mode)

### 1. Install dependencies
```bash
go mod tidy
```

### 2. Generate Swagger docs
```bash
make swag
```
(or manually: `swag init -g cmd/server/main.go -o ./docs`)

### 3. Run server
```bash
make run
```
(or manually: `go run ./cmd/server`)

Server will start on:
```
http://localhost:8080
```

---

## 📑 Swagger Docs
Open browser:
```
http://localhost:8080/swagger/index.html
```

Or fetch raw spec:
```bash
curl http://localhost:8080/swagger/doc.json | jq .
```

---

## 🔑 API Usage Examples

### 1. Request OTP
```bash
curl -X POST http://localhost:8080/api/v1/auth/request-otp \
  -H 'Content-Type: application/json'   -d '{"phone":"+15551234567"}'
```
👉 Response:
```json
{"message": "otp generated (check server logs)"}
```
👉 OTP is printed in **server console logs**, e.g.:
```
[OTP] +15551234567 -> 123456 (expires in 2m0s)
```

---

### 2. Verify OTP (login/register → JWT)
```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H 'Content-Type: application/json'   -d '{"phone":"+15551234567","otp":"123456"}'
```
👉 Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "c2b45c0e-...",
    "phone": "+15551234567",
    "registered_at": "2025-08-30T20:00:00Z"
  }
}
```

---

### 3. Access Protected Endpoints (requires JWT)

Get token from previous step:

```bash
TOKEN=<your_jwt>
```

#### a) Get User by ID
```bash
curl -H "Authorization: Bearer $TOKEN"   http://localhost:8080/api/v1/users/<user_id>
```

#### b) List Users
```bash
curl -H "Authorization: Bearer $TOKEN"   "http://localhost:8080/api/v1/users?page=1&size=10&search=+1555"
```

👉 Example response:
```json
{
  "items": [
    {
      "id": "c2b45c0e-...",
      "phone": "+15551234567",
      "registered_at": "2025-08-30T20:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 10
}
```

---

## 🐳 Run with Docker (in-memory mode)

### 1. Build image
```bash
docker build -t otp-service:dev .
```

### 2. Run container
```bash
docker run --rm -p 8080:8080 otp-service:dev
```

Server inside container is now available at:
```
http://localhost:8080/swagger/index.html
```

---

## ⚠️ Notes / Limitations

- **Data persistence**: In-memory repo means all users are lost when service restarts.  
- **Scaling**: OTPs & rate limits are stored in-memory, so not suitable for multi-instance deployments.  
- **Security**: OTPs are printed to logs (demo purpose). In production integrate with SMS/email provider.  
- **Secrets**: Always set a strong `JWT_SECRET` via environment variable.

---

## 🔜 Next Steps
- Add **PostgreSQL support** with `docker-compose`
- Persist users in DB
- Store OTPs & rate-limits in **Redis** for distributed deployments
