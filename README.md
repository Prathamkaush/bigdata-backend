# 🚀 BigData API — High-Performance Credit-Based Big Data Query System  
A scalable backend that supports **big-data ingestion**, **ClickHouse analytics**, **PostgreSQL user/credit system**, **Redis caching & rate-limiting**, and a secure **credit-deducting API layer**.

This project powers the entire big-data admin system allowing:
- Multi-format ingestion (CSV / JSON / XML / Parquet)
- Millions of records stored in ClickHouse
- Credit-based API access for users
- Admin controls: users, credits, logs
- Full logging + analytics dashboard for the frontend

---

## 🏗️ **Tech Stack**
### **Backend**
- **Go (Golang)**
- **Fiber v2**
- **PostgreSQL (NeonDB)** → users, credits, logs
- **ClickHouse Cloud** → big data & analytics
- **Upstash Redis** → rate limiting + caching
- **JWT-less API key authorization**
- Fully layered architecture: `controllers → services → repositories → database`

---

## 📁 **Project Structure**

```
/bigdata-api
│
├── cmd/server/main.go
├── internal/
│   ├── config/           → Loads environment variables
│   ├── database/         → Postgres / Redis / ClickHouse connectors
│   ├── api/
│   │   ├── middlewares/  → Auth, admin, credits, logging, rate-limit
│   │   ├── controllers/  → Admin, Stats, Logs, Query
│   │   └── routes/       → Route definitions
│   ├── services/         → User, credit, query service
│   ├── repository/       → DB operations
│   ├── ingestion/        → CSV/JSON/XML ingestion pipeline
│   ├── models/           → Structs for DB + API
│   └── utils/            → hasher, response formatter
│
├── scripts/              → DB migration & ingestion scripts
├── .env
├── go.mod
└── go.sum
```

---

## 🔑 **Environment Variables (.env)**

```
SERVER_PORT=8080

# PostgreSQL
POSTGRES_URL=postgresql://USER:PASSWORD@HOST/neondb?sslmode=require

# ClickHouse
CLICKHOUSE_HOST=xxxx.ap-south-1.aws.clickhouse.cloud:8443
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=xxxxxx

# Redis
REDIS_URL=rediss://default:xxxxx@xxx.upstash.io:6379

# Admin API Key (SHA256 hash)
ADMIN_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

## 🧪 **Run Locally**

### 1️⃣ Install dependencies
```bash
go mod tidy
```

### 2️⃣ Run server
```bash
go run ./cmd/server
```

Server runs at:  
👉 **http://localhost:8080**

---

## 📡 **API Endpoints**

### **Admin**
```
POST   /v1/admin/create-user
POST   /v1/admin/add-credits
GET    /v1/admin/users
GET    /v1/admin/logs
GET    /v1/admin/stats
```

### **User Query**
```
POST   /v1/query
```

### **Required Headers**
```
x-api-key: USER_API_KEY
```

---

## 💳 **Credit System**
| Event | Credits |
|-------|---------|
| Query API request | -1 credit |
| Credits reach 0 | API returns "Insufficient Credits" |
| Admin can recharge credits | ✔ |

---

## 📈 Logging & Analytics
Every request logs:
- endpoint  
- user_id  
- timestamp  
- status  
- duration (ms)  

Stored in **Postgres** → used in frontend dashboard.

---

## 🌐 Deployment (Render)
1. Add environment variables in Render Dashboard  
2. Set build command:
```bash
go build -o app ./cmd/server
```
3. Start command:
```bash
./app
```

---

## 🧑‍💻 Author
**Pratham Kaushik**  
Big Data API Architect & Full Stack Developer  
GitHub: https://github.com/Prathamkaush

