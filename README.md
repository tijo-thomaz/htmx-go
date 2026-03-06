# LinkBio — Linktree Clone

A **production-grade** link-in-bio platform built with **Go**, **HTMX**, **Alpine.js**, **AOS**, and **GSAP**.

> 🎬 **YouTube Tutorial Project** — Build this complete app in 1 hour (Malayalam)

---

## 🏭 Production-Grade Go Patterns

This tutorial teaches **industry-level** Go practices:

| Pattern | Description |
|---------|-------------|
| **Namespace-based routing** | Organized route groups (`/api/v1`, `/auth`, `/admin`) |
| **Structured logging** | Log levels (DEBUG, INFO, WARN, ERROR) with `slog` |
| **Graceful shutdown** | Handle SIGINT/SIGTERM, drain connections |
| **Context propagation** | Request-scoped values, timeouts, cancellation |
| **Connection pooling** | SQLite connection pool for concurrent users |
| **Rate limiting** | Protect against abuse |
| **Middleware chain** | Logger → Recovery → Auth → Handler |
| **Error handling** | Centralized error responses |
| **Config management** | Environment-based configuration |
| **Health checks** | `/health` endpoint for monitoring |

---

## 🎯 Project Overview

LinkBio allows users to create a single page with all their important links. Share one URL everywhere — Instagram, YouTube, TikTok, WhatsApp — and let your audience find everything in one place.

**Live Demo:** `yourdomain.com/u/username`

---

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go** | Backend server, auth, database, API |
| **HTMX** | Dynamic updates without full page reload |
| **Alpine.js** | Client-side interactivity (drag, toggle, modals) |
| **AOS** | Scroll-based animations on public page |
| **GSAP** | Smooth hover effects and page transitions |
| **Tailwind CSS** | Styling |
| **SQLite** | Database |

---

## ✨ Features

### Core Features (Tutorial Scope)

| Feature | Description | Tech Used |
|---------|-------------|-----------|
| **User Auth** | Signup, login, logout with sessions | Go |
| **Add Links** | Create new links without page reload | HTMX |
| **Edit Links** | Inline edit link title and URL | HTMX |
| **Delete Links** | Remove links with confirmation | HTMX |
| **Drag Reorder** | Reorder links by dragging | Alpine.js |
| **Theme Toggle** | Light/Dark mode switch | Alpine.js |
| **Public Profile** | Shareable page at `/u/username` | Go Templates |
| **Link Animations** | Fade-in on scroll | AOS |
| **Hover Effects** | Smooth button animations | GSAP |
| **Click Analytics** | Track link clicks and page views | Go |
| **Deploy** | Live deployment | Railway/Fly.io |

### Future Enhancements (Not in Tutorial)

- [ ] Google OAuth login
- [ ] Stripe payments & pricing tiers
- [ ] Email collection
- [ ] Custom domains
- [ ] Video embeds
- [ ] Scheduled links
- [ ] Advanced analytics (location, device, referrer)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         BROWSER                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Alpine.js (client state, drag-reorder, theme)      │    │
│  │  GSAP (animations)                                   │    │
│  │  AOS (scroll reveals)                                │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTMX (HTML over the wire)
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       GO SERVER                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Routes  │  │ Handlers │  │ Sessions │  │ Templates│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       SQLite                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │  users   │  │  links   │  │ analytics│                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure (Industry Standard)

```
linkbio/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── server/
│   │   └── server.go            # HTTP server setup
│   ├── router/
│   │   ├── router.go            # Main router setup
│   │   ├── auth.go              # /auth routes namespace
│   │   ├── api.go               # /api/v1 routes namespace
│   │   └── public.go            # Public routes namespace
│   ├── handler/
│   │   ├── auth.go              # Auth handlers
│   │   ├── link.go              # Link CRUD handlers
│   │   ├── profile.go           # Public profile handler
│   │   └── health.go            # Health check handler
│   ├── middleware/
│   │   ├── logger.go            # Request logging middleware
│   │   ├── recovery.go          # Panic recovery middleware
│   │   ├── auth.go              # Auth middleware
│   │   └── ratelimit.go         # Rate limiting middleware
│   ├── model/
│   │   ├── user.go              # User struct & methods
│   │   ├── link.go              # Link struct & methods
│   │   └── analytics.go         # Analytics struct & methods
│   ├── repository/
│   │   ├── user.go              # User database operations
│   │   ├── link.go              # Link database operations
│   │   └── analytics.go         # Analytics database operations
│   ├── service/
│   │   ├── auth.go              # Auth business logic
│   │   ├── link.go              # Link business logic
│   │   └── analytics.go         # Analytics business logic
│   └── pkg/
│       ├── logger/
│       │   └── logger.go        # Structured logging (slog)
│       ├── response/
│       │   └── response.go      # Standardized HTTP responses
│       └── validator/
│           └── validator.go     # Input validation
├── web/
│   ├── templates/
│   │   ├── layouts/
│   │   │   └── base.html        # Base layout
│   │   ├── pages/
│   │   │   ├── home.html        # Landing page
│   │   │   ├── login.html       # Login form
│   │   │   ├── register.html    # Signup form
│   │   │   ├── dashboard.html   # Admin panel
│   │   │   └── profile.html     # Public profile
│   │   └── partials/
│   │       ├── link.html        # Single link (HTMX)
│   │       ├── link-form.html   # Add/Edit form
│   │       └── analytics.html   # Stats display
│   └── static/
│       ├── css/
│       │   └── style.css        # Tailwind CSS
│       ├── js/
│       │   └── app.js           # Alpine + GSAP + AOS
│       └── images/
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_links.sql
│   └── 003_create_analytics.sql
├── .env.example
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 🗄️ Database Schema

### users
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name TEXT,
    bio TEXT,
    avatar_url TEXT,
    theme TEXT DEFAULT 'light',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### links
```sql
CREATE TABLE links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    icon TEXT,
    position INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### analytics
```sql
CREATE TABLE analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    link_id INTEGER,
    event_type TEXT NOT NULL,  -- 'page_view' or 'link_click'
    referrer TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (link_id) REFERENCES links(id)
);
```

---

## 🔧 Production Go Patterns (Code Examples)

### Graceful Shutdown

```go
// cmd/server/main.go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    srv := server.New(cfg)
    
    // Start server in goroutine
    go func() {
        log.Info("server starting", "port", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Error("server error", "error", err)
        }
    }()

    // Wait for interrupt signal
    <-ctx.Done()
    log.Info("shutdown signal received")

    // Graceful shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Error("shutdown error", "error", err)
    }
    log.Info("server stopped gracefully")
}
```

### Structured Logging (slog)

```go
// internal/pkg/logger/logger.go
package logger

import (
    "log/slog"
    "os"
)

type Level string

const (
    DEBUG Level = "DEBUG"
    INFO  Level = "INFO"
    WARN  Level = "WARN"
    ERROR Level = "ERROR"
)

func New(level Level) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case DEBUG:
        logLevel = slog.LevelDebug
    case INFO:
        logLevel = slog.LevelInfo
    case WARN:
        logLevel = slog.LevelWarn
    case ERROR:
        logLevel = slog.LevelError
    }

    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    })
    return slog.New(handler)
}

// Usage:
// log.Info("user created", "user_id", user.ID, "email", user.Email)
// log.Error("database error", "error", err, "query", query)
```

### Namespace-based Routing

```go
// internal/router/router.go
func New(h *handler.Handler, mw *middleware.Middleware) *chi.Mux {
    r := chi.NewRouter()

    // Global middleware
    r.Use(mw.Logger)
    r.Use(mw.Recovery)
    r.Use(mw.RateLimit)

    // Health check (no auth)
    r.Get("/health", h.Health.Check)

    // Public routes
    r.Group(func(r chi.Router) {
        r.Get("/", h.Home.Index)
        r.Get("/u/{username}", h.Profile.Show)
    })

    // Auth namespace
    r.Route("/auth", func(r chi.Router) {
        r.Get("/login", h.Auth.LoginPage)
        r.Post("/login", h.Auth.Login)
        r.Get("/register", h.Auth.RegisterPage)
        r.Post("/register", h.Auth.Register)
        r.Post("/logout", h.Auth.Logout)
    })

    // API v1 namespace (protected)
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(mw.Auth) // Require authentication
        
        r.Route("/links", func(r chi.Router) {
            r.Post("/", h.Link.Create)
            r.Put("/{id}", h.Link.Update)
            r.Delete("/{id}", h.Link.Delete)
            r.Post("/reorder", h.Link.Reorder)
        })

        r.Get("/analytics", h.Analytics.Dashboard)
    })

    // Admin namespace
    r.Route("/dashboard", func(r chi.Router) {
        r.Use(mw.Auth)
        r.Get("/", h.Dashboard.Index)
    })

    return r
}
```

### Middleware Chain

```go
// internal/middleware/logger.go
func (m *Middleware) Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status
        ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
        
        next.ServeHTTP(ww, r)
        
        m.log.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "status", ww.status,
            "duration", time.Since(start),
            "ip", r.RemoteAddr,
        )
    })
}

// internal/middleware/recovery.go
func (m *Middleware) Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                m.log.Error("panic recovered",
                    "error", err,
                    "stack", string(debug.Stack()),
                )
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Connection Pool & Concurrent Users

```go
// internal/repository/db.go
func NewDB(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
    if err != nil {
        return nil, err
    }

    // Connection pool settings for concurrent users
    db.SetMaxOpenConns(25)              // Max open connections
    db.SetMaxIdleConns(10)              // Max idle connections
    db.SetConnMaxLifetime(5 * time.Minute) // Connection max lifetime
    db.SetConnMaxIdleTime(1 * time.Minute) // Idle connection timeout

    // Enable WAL mode for better concurrent read/write
    if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
        return nil, err
    }

    return db, nil
}
```

### Context Propagation

```go
// internal/handler/link.go
func (h *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Get user from context (set by auth middleware)
    user := middleware.UserFromContext(ctx)
    if user == nil {
        h.response.Error(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    // Context with timeout for database operations
    dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    link, err := h.service.Create(dbCtx, user.ID, req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            h.response.Error(w, http.StatusGatewayTimeout, "database timeout")
            return
        }
        h.response.Error(w, http.StatusInternalServerError, "failed to create link")
        return
    }

    h.response.HTML(w, http.StatusCreated, "partials/link.html", link)
}
```

### Rate Limiting

```go
// internal/middleware/ratelimit.go
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
    // Simple in-memory rate limiter (use Redis for production)
    limiter := rate.NewLimiter(rate.Limit(10), 20) // 10 req/sec, burst 20

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            m.log.Warn("rate limit exceeded", "ip", r.RemoteAddr)
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Config Management

```go
// internal/config/config.go
type Config struct {
    Port          string
    Env           string // development, staging, production
    LogLevel      string
    DatabasePath  string
    SessionSecret string
    RateLimit     int
}

func Load() (*Config, error) {
    return &Config{
        Port:          getEnv("PORT", "8080"),
        Env:           getEnv("ENV", "development"),
        LogLevel:      getEnv("LOG_LEVEL", "INFO"),
        DatabasePath:  getEnv("DATABASE_PATH", "./linkbio.db"),
        SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
        RateLimit:     getEnvInt("RATE_LIMIT", 10),
    }, nil
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
```

---

## 🛣️ Routes (Namespace-based)

### Public Routes
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `Home.Index` | Landing page |
| GET | `/u/{username}` | `Profile.Show` | Public profile page |
| GET | `/health` | `Health.Check` | Health check endpoint |

### Auth Namespace (`/auth`)
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/auth/login` | `Auth.LoginPage` | Login form |
| POST | `/auth/login` | `Auth.Login` | Authenticate user |
| GET | `/auth/register` | `Auth.RegisterPage` | Signup form |
| POST | `/auth/register` | `Auth.Register` | Create account |
| POST | `/auth/logout` | `Auth.Logout` | End session |

### API Namespace (`/api/v1`) — Protected
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/api/v1/links` | `Link.Create` | Add new link (HTMX) |
| PUT | `/api/v1/links/{id}` | `Link.Update` | Edit link (HTMX) |
| DELETE | `/api/v1/links/{id}` | `Link.Delete` | Remove link (HTMX) |
| POST | `/api/v1/links/reorder` | `Link.Reorder` | Update positions |
| POST | `/api/v1/analytics/click/{id}` | `Analytics.Click` | Track link click |
| GET | `/api/v1/analytics` | `Analytics.Dashboard` | Get stats |
| GET | `/api/v1/qr` | `QR.Generate` | Generate QR code |

### Dashboard Namespace (`/dashboard`) — Protected
| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/dashboard` | `Dashboard.Index` | Admin panel |

---

## 🎨 UI/UX Specifications

### Color Palette

```css
/* Light Theme */
--bg-primary: #ffffff;
--bg-secondary: #f3f4f6;
--text-primary: #111827;
--text-secondary: #6b7280;
--accent: #3b82f6;

/* Dark Theme */
--bg-primary: #0f172a;
--bg-secondary: #1e293b;
--text-primary: #f8fafc;
--text-secondary: #94a3b8;
--accent: #22d3ee;
```

### Typography

- **Headings:** Inter or system-ui
- **Body:** Inter or system-ui
- **Code (if any):** JetBrains Mono

### Animations

| Element | Animation | Library |
|---------|-----------|---------|
| Links on public page | Fade up on scroll | AOS |
| Button hover | Scale + glow | GSAP |
| Theme switch | Smooth transition | CSS + Alpine |
| Link add/delete | Slide in/out | HTMX swap |
| Page load | Staggered fade | GSAP |

---

## 📺 Tutorial Timeline (1 Hour)

| Time | Section | Content | Go Patterns Covered |
|------|---------|---------|---------------------|
| 0:00-3:00 | **Hook** | Show finished app demo | — |
| 3:00-8:00 | **Architecture** | Draw diagram, explain tech stack | Namespace routing concept |
| 8:00-18:00 | **Setup** | Project structure, config, logger, server | Config management, Structured logging, Graceful shutdown |
| 18:00-28:00 | **Auth** | Signup, login, sessions, middleware | Middleware chain, Context propagation |
| 28:00-38:00 | **HTMX Links** | Add, edit, delete without reload | Repository pattern, Error handling |
| 38:00-45:00 | **Alpine.js** | Drag reorder, theme toggle | — |
| 45:00-50:00 | **Public Page** | Profile template, AOS, GSAP | Rate limiting, Connection pooling |
| 50:00-55:00 | **Analytics + QR** | Click tracking, QR generation | Concurrent writes |
| 55:00-58:00 | **Deploy** | Railway/Fly.io deployment | Health checks |
| 58:00-60:00 | **Recap** | Summary, next steps, CTA | — |

### What Viewers Will Learn

```
✅ Industry-standard Go project structure
✅ Graceful shutdown (handle Ctrl+C properly)
✅ Structured logging with log levels (slog)
✅ Namespace-based routing (/auth, /api/v1, /dashboard)
✅ Middleware chain (Logger → Recovery → RateLimit → Auth)
✅ Context propagation (timeouts, cancellation)
✅ Connection pooling for concurrent users
✅ Rate limiting to prevent abuse
✅ Health check endpoints for monitoring
✅ Environment-based configuration
✅ HTMX for dynamic updates without JavaScript
✅ Alpine.js for client-side interactivity
✅ AOS + GSAP for professional animations
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.22+
- Make (optional, for Makefile commands)

### Installation

```bash
# Clone the repo
git clone https://github.com/tijo-thomaz/htmx-go.git
cd linkbio

# Install Go dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Run the server
make run
# OR
go run ./cmd/server

# Open in browser
# http://localhost:8080
```

### Environment Variables

```env
# .env.example
PORT=8080
ENV=development
LOG_LEVEL=DEBUG
DATABASE_PATH=./data/linkbio.db
SESSION_SECRET=your-super-secret-key-change-in-production
RATE_LIMIT=10
```

### Makefile Commands

```makefile
# Makefile
.PHONY: run build test clean

run:
	go run ./cmd/server

build:
	go build -o bin/linkbio.exe ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/ data/

dev:
	air  # Hot reload (requires air installed)

migrate:
	go run ./cmd/migrate
```

---

## 📦 Dependencies

```go
// go.mod
module linkbio

go 1.22

require (
    github.com/go-chi/chi/v5 v5.0.12       // Router with namespace support
    github.com/gorilla/sessions v1.2.2     // Session management
    modernc.org/sqlite v1.29.5             // SQLite driver (pure Go, no CGO)
    golang.org/x/crypto v0.21.0            // Password hashing (bcrypt)
    golang.org/x/time v0.5.0               // Rate limiting
    github.com/joho/godotenv v1.5.1        // Environment file loading
)
```

### Frontend (CDN)

```html
<!-- HTMX -->
<script src="https://unpkg.com/htmx.org@1.9.6"></script>

<!-- Alpine.js -->
<script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>

<!-- AOS -->
<link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet">
<script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>

<!-- GSAP -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/gsap/3.12.2/gsap.min.js"></script>

<!-- Tailwind (CDN for dev) -->
<script src="https://cdn.tailwindcss.com"></script>
```

---

## 🎬 YouTube Channel Info

**Channel Name:** [Your Channel Name]  
**Language:** Malayalam  
**Style:** Faceless, code-focused, diagram-first explanations  
**Target:** Beginners to Senior developers  

### Teaching Approach

1. **Diagram First** — Draw architecture before coding
2. **Why Before How** — Explain the problem, then solution
3. **Triple Explanation** — Technical → Simple → Analogy
4. **Malayalam + English** — Malayalam for concepts, English for code

---

## 📄 License

MIT License — Free to use, modify, and distribute.

---

## 🤝 Contributing

This is a tutorial project. Feel free to:
- Fork and extend
- Report issues
- Suggest improvements

---

---

## 🧪 Testing the Production Patterns

### Test Graceful Shutdown
```bash
# Start server
make run

# In another terminal, send SIGTERM
kill -SIGTERM $(pgrep linkbio)

# Or press Ctrl+C - watch logs for graceful shutdown
```

### Test Rate Limiting
```bash
# Send 25 requests quickly
for i in {1..25}; do curl -s http://localhost:8080/health; done

# Should see "Too Many Requests" after limit exceeded
```

### Test Concurrent Users
```bash
# Install hey (HTTP load generator)
go install github.com/rakyll/hey@latest

# Send 100 concurrent requests
hey -n 1000 -c 100 http://localhost:8080/u/testuser
```

### Check Logs
```bash
# Logs are JSON formatted for production
go run ./cmd/server 2>&1 | jq .

# Sample output:
# {"time":"2024-...","level":"INFO","msg":"request","method":"GET","path":"/health","status":200,"duration":"1.2ms"}
```

---

## 📚 Additional Resources

- [HTMX Documentation](https://htmx.org/docs/)
- [Go Chi Router](https://github.com/go-chi/chi)
- [Alpine.js Documentation](https://alpinejs.dev/)
- [AOS - Animate on Scroll](https://michalsnik.github.io/aos/)
- [GSAP Animation](https://greensock.com/gsap/)
- [Go slog (Structured Logging)](https://pkg.go.dev/log/slog)

---

Built with ❤️ for the Malayalam developer community.
