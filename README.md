# Instagram Clone Architecture & Implementation

An Instagram-scale full-stack architecture powered by:
- **Core Security & Business Domain**: Java 21 + Spring Boot 3.3 + Spring Data MongoDB + Spring Security 6 (JWT)
- **High-Concurrency Real-Time Engine**: Go (Gorilla WebSockets, lightweight connection pooling, direct message routing)
- **Database**: MongoDB (Local or MongoDB Atlas)
- **Frontend**: Clean Instagram registration UI with smooth CSS transitions and live API integration

---

## Project Structure

```
├── backend-core/               # Java 21 & Spring Boot Core Service
│   ├── pom.xml
│   └── src/main/java/com/instagram/core/
│       ├── controller/        # AuthController (/api/v1/auth/*), UserController
│       ├── dto/               # RegisterRequest, LoginRequest, AuthResponse, UserProfileDto
│       ├── model/             # User Document, Role Enum
│       ├── repository/        # UserRepository (MongoDB)
│       ├── security/          # JwtService, JwtAuthenticationFilter, SecurityConfig
│       └── service/           # AuthService, UserService
│
├── backend-concurrency/        # Go Real-time Concurrency Service
│   ├── main.go                # HTTP & WebSocket endpoints (/ws, /health, /api/v1/stats, /api/v1/notify)
│   ├── go.mod / go.sum
│   └── hub/
│       ├── hub.go             # Concurrency Hub with thread-safe client tracking
│       └── client.go          # WebSocket read/write pumps with heartbeat
│
├── index.html                 # Instagram Web Registration Page with smooth CSS transitions
└── README.md
```

---

## How to Run

### 1. Prerequisites
- **MongoDB**: Already running on `localhost:27017` (Database: `instagram_db`)
- **Java**: Java 21 LTS
- **Maven**: 3.9+
- **Go**: 1.22+

---

### 2. Run the Core Service (Java 21 / Spring Boot)
```powershell
cd backend-core
mvn spring-boot:run
```
- Core service will start on **http://localhost:8080**
- API endpoints:
  - `POST /api/v1/auth/register` - Create account & issue JWT
  - `POST /api/v1/auth/login` - Authenticate & issue JWT
  - `GET  /api/v1/auth/me` - Get current logged-in user (requires `Authorization: Bearer <token>`)
  - `GET  /api/v1/users/{username}` - Get public user profile

---

### 3. Run the High-Concurrency Service (Go)
```powershell
cd backend-concurrency
.\concurrency-service.exe
# OR: go run .
```
- Real-time service will start on **http://localhost:8081**
- Endpoints:
  - `GET  /health` - Health check
  - `GET  /ws?user_id=<user_id>` - WebSocket connection for live DMs and notifications
  - `GET  /api/v1/stats` - Live active connections and goroutines
  - `POST /api/v1/notify` - Internal push event dispatcher

---

### 4. Test the Frontend
Open [index.html](file:///D:/Desktop/gofiles/faltu/index.html) in your browser.
- Fill out Mobile/Email, Full Name, Username, and Password.
- Clicking **Sign up** will automatically call the Spring Boot registration endpoint and store the JWT in your browser's `localStorage`.
