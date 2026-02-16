# LinkBio Tutorial - Complete Narrator Script

> 🎬 **YouTube Tutorial**: Build a Production-Grade Link-in-Bio Platform
> ⏱️ **Duration**: ~60 minutes
> 🗣️ **Language**: Malayalam (Code in English)

---

## 📋 Pre-Recording Checklist

- [ ] VS Code open with empty folder
- [ ] Terminal ready
- [ ] Go 1.22+ installed (`go version`)
- [ ] Browser open for testing
- [ ] This script on phone

---

# PART 1: INTRODUCTION (0:00 - 3:00)

## Opening Hook

> "നമസ്കാരം! ഇന്ന് നമ്മൾ ഒരു production-ready link-in-bio platform build ചെയ്യാൻ പോകുന്നു - Linktree പോലെ. ഇത് ഒരു beginner project അല്ല - ഇത് real companies use ചെയ്യുന്ന patterns ആണ്."

**Show the finished app demo:**
- Landing page with animations
- User registration/login
- Dashboard with links
- Public profile page
- Dark mode toggle

> "ഈ tutorial-ൽ നിങ്ങൾ പഠിക്കുന്നത്:
> - Go backend with industry patterns
> - HTMX for dynamic updates - JavaScript എഴുതാതെ!
> - Alpine.js for client-side magic
> - Beautiful UI with Tailwind CSS
> - Real database with SQLite"

---

# PART 2: PROJECT ARCHITECTURE (3:00 - 8:00)

## Draw This Diagram First

```
┌─────────────────────────────────────────────────────────────┐
│                         BROWSER                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Alpine.js (toggle, drag)  │  GSAP (animations)      │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │ HTMX (sends HTML, not JSON!)
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       GO SERVER                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Router  │→ │Middleware│→ │ Handlers │→ │ Templates│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                       SQLite Database                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │  users   │  │  links   │  │ analytics│                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

## Explain the Architecture

> "ഈ diagram നോക്കൂ. നമ്മുടെ app-ന് മൂന്ന് layers ഉണ്ട്:"

**Layer 1 - Browser:**
> "User-ന്റെ browser. ഇവിടെ Alpine.js toggle, drag-drop handle ചെയ്യും. GSAP animations കാണിക്കും."

**Layer 2 - Go Server:**
> "ഇത് നമ്മുടെ brain. Request വരുമ്പോൾ Router ആദ്യം handle ചെയ്യും. പിന്നെ Middleware - ഇത് security guard പോലെ. പിന്നെ Handler - actual work. അവസാനം Template - HTML generate."

**Layer 3 - Database:**
> "SQLite - ഒരു file-based database. Production-ന് PostgreSQL use ചെയ്യാം, പക്ഷേ learning-ന് SQLite best."

---

# PART 3: PROJECT SETUP (8:00 - 15:00)

## Step 1: Create Project Folder

**📱 Read this:**
> "ആദ്യം നമ്മൾ project folder create ചെയ്യാം."

**⌨️ Type this:**
```bash
mkdir linkbio
cd linkbio
go mod init linkbio
```

**🧠 Explain:**
> "go mod init - ഇത് നമ്മുടെ project-ന്റെ ID card പോലെ. Package management ഇത് handle ചെയ്യും."

---

## Step 2: Create Folder Structure

**📱 Read this:**
> "Industry-standard folder structure create ചെയ്യാം. ഇത് large companies use ചെയ്യുന്ന pattern ആണ്."

**⌨️ Type this (one by one):**
```bash
mkdir -p cmd/server
mkdir -p internal/config
mkdir -p internal/server
mkdir -p internal/router
mkdir -p internal/handler
mkdir -p internal/middleware
mkdir -p internal/model
mkdir -p internal/repository
mkdir -p internal/pkg/logger
mkdir -p internal/pkg/response
mkdir -p internal/pkg/templates
mkdir -p web/templates/layouts
mkdir -p web/templates/pages
mkdir -p web/templates/partials
mkdir -p web/static/css
mkdir -p web/static/js
mkdir -p data
```

**🧠 Explain:**
> "cmd/ - entry points. നമ്മുടെ main.go ഇവിടെ."
> "internal/ - private code. മറ്റ് projects import ചെയ്യാൻ പറ്റില്ല."
> "web/ - HTML, CSS, JS files."
> "data/ - SQLite database file."

**🎯 Why this structure?**
> "ഇത് random അല്ല. Google, Uber, Stripe - എല്ലാവരും similar structure use ചെയ്യുന്നു. ഇത് scalable, maintainable."

---

## Step 3: Install Dependencies

**📱 Read this:**
> "നമുക്ക് ആവശ്യമായ packages install ചെയ്യാം."

**⌨️ Type this:**
```bash
go get github.com/go-chi/chi/v5
go get github.com/gorilla/sessions
go get modernc.org/sqlite
go get golang.org/x/crypto
go get golang.org/x/time
go get github.com/joho/godotenv
```

**🧠 Explain each:**
> "chi - Router. Express.js പോലെ, പക്ഷേ Go-യ്ക്ക്."
> "gorilla/sessions - Login sessions handle ചെയ്യാൻ."
> "modernc/sqlite - Database. Pure Go, no C compiler needed."
> "x/crypto - Password hashing. bcrypt."
> "x/time - Rate limiting."
> "godotenv - .env file reading."

---

## Step 4: Create .env File

**📱 Read this:**
> "Configuration file create ചെയ്യാം. Secrets ഇവിടെ store ചെയ്യും."

**⌨️ Create `.env`:**
```env
PORT=8080
ENV=development
LOG_LEVEL=DEBUG
DATABASE_PATH=./data/linkbio.db
SESSION_SECRET=super-secret-key-change-in-production-minimum-32-chars
RATE_LIMIT=100
```

**🧠 Explain:**
> "ഇത് environment variables. Production-ൽ different values ആകും. Code-ൽ hardcode ചെയ്യരുത് - security risk!"

**⚠️ Important:**
> "SESSION_SECRET minimum 32 characters വേണം. Production-ൽ random generate ചെയ്യണം."

---

# PART 4: CONFIGURATION (15:00 - 18:00)

## Create config.go

**📱 Read this:**
> "Configuration loading code എഴുതാം. ഇത് .env file read ചെയ്യും."

**⌨️ Create `internal/config/config.go`:**
```go
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	Port          string
	Env           string
	LogLevel      string
	DatabasePath  string
	SessionSecret string
	RateLimit     int
}

// Load reads configuration from environment
func Load() (*Config, error) {
	// Load .env file (ignore error if not exists)
	godotenv.Load()

	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT", "100"))

	return &Config{
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "INFO"),
		DatabasePath:  getEnv("DATABASE_PATH", "./data/linkbio.db"),
		SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
		RateLimit:     rateLimit,
	}, nil
}

// getEnv reads env variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// IsDevelopment returns true if running in dev mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}
```

**🧠 Explain:**
> "getEnv function - environment variable read ചെയ്യും. ഇല്ലെങ്കിൽ default value return ചെയ്യും."
> "Config struct - എല്ലാ settings ഒരു place-ൽ."
> "IsDevelopment - dev/prod check ചെയ്യാൻ."

---

# PART 5: LOGGING (18:00 - 20:00)

## Why Structured Logging?

**📱 Read this:**
> "Production apps-ൽ fmt.Println use ചെയ്യരുത്! Structured logging വേണം."

**🎯 Analogy:**
> "fmt.Println ഒരു diary പോലെ - random notes. Structured logging ഒരു spreadsheet പോലെ - organized, searchable."

## Create logger.go

**⌨️ Create `internal/pkg/logger/logger.go`:**
```go
package logger

import (
	"log/slog"
	"os"
)

// New creates a production logger
func New(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

// NewDevelopment creates a dev-friendly logger
func NewDevelopment() *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
```

**🧠 Explain:**
> "slog - Go 1.21-ൽ വന്ന built-in structured logger."
> "Production-ൽ JSON format - log aggregators (Datadog, Splunk) parse ചെയ്യാൻ."
> "Development-ൽ Text format - human readable."

---

# PART 6: DATABASE (20:00 - 25:00)

## Create Models

**📱 Read this:**
> "ആദ്യം data structures define ചെയ്യാം. ഇത് database tables-ന്റെ Go representation ആണ്."

**⌨️ Create `internal/model/user.go`:**
```go
package model

import "time"

// User represents a user account
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send to client!
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	Theme        string    `json:"theme"`
	CreatedAt    time.Time `json:"created_at"`
}
```

**🧠 Explain:**
> "json:\"-\" - PasswordHash JSON-ൽ include ചെയ്യരുത്. Security!"

**⌨️ Create `internal/model/link.go`:**
```go
package model

import "time"

// Link represents a user's link
type Link struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Icon      string    `json:"icon"`
	Position  int       `json:"position"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
```

**⌨️ Create `internal/model/analytics.go`:**
```go
package model

import "time"

// Analytics represents a tracking event
type Analytics struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LinkID    *int64    `json:"link_id"` // Pointer because nullable
	EventType string    `json:"event_type"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsSummary for dashboard
type AnalyticsSummary struct {
	TotalViews  int `json:"total_views"`
	TotalClicks int `json:"total_clicks"`
}
```

---

## Create Database Connection

**📱 Read this:**
> "Database connection code എഴുതാം. Connection pooling, migrations - എല്ലാം handle ചെയ്യും."

**⌨️ Create `internal/repository/db.go`:**
```go
package repository

import (
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

// NewDB creates a database connection with pooling
func NewDB(path string, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info("database connected", "path", path)
	return db, nil
}

// Migrate creates tables
func Migrate(db *sql.DB, log *slog.Logger) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			display_name TEXT DEFAULT '',
			bio TEXT DEFAULT '',
			avatar_url TEXT DEFAULT '',
			theme TEXT DEFAULT 'light',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			icon TEXT DEFAULT '',
			position INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS analytics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			link_id INTEGER,
			event_type TEXT NOT NULL,
			referrer TEXT DEFAULT '',
			user_agent TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_links_user_id ON links(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_user_id ON analytics(user_id)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Error("migration failed", "index", i, "error", err)
			return err
		}
	}

	log.Info("database migrations completed")
	return nil
}
```

**🧠 Explain key concepts:**

> "WAL mode - Write-Ahead Logging. Multiple readers, one writer. Faster!"

> "Connection pooling - ഓരോ request-നും new connection വേണ്ട. Pool-ൽ നിന്ന് reuse ചെയ്യും. Restaurant-ൽ plates wash and reuse ചെയ്യുന്നത് പോലെ."

> "Migrations - Tables create ചെയ്യുന്ന code. App start ചെയ്യുമ്പോൾ auto-run ആകും."

---

## Create User Repository

**⌨️ Create `internal/repository/user.go`:**
```go
package repository

import (
	"context"
	"database/sql"

	"linkbio/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, display_name, theme)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Theme,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// GetByID finds user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, theme, created_at FROM users WHERE id = ?`
	
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.Theme, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByUsername finds user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, theme, created_at FROM users WHERE username = ?`
	
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.Theme, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail finds user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, display_name, bio, avatar_url, theme, created_at FROM users WHERE email = ?`
	
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.Theme, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
```

**🧠 Explain:**
> "Repository pattern - Database code ഒരു place-ൽ. Handler-ൽ SQL എഴുതരുത്. Separation of concerns."

> "Context - Request cancel ആയാൽ database query-ഉം cancel ആകും. User browser close ചെയ്താൽ server resources waste ചെയ്യില്ല."

---

## Create Link Repository

**⌨️ Create `internal/repository/link.go`:**
```go
package repository

import (
	"context"
	"database/sql"

	"linkbio/internal/model"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// Create inserts a new link
func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	// Get next position
	var maxPos int
	r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(position), 0) FROM links WHERE user_id = ?", link.UserID).Scan(&maxPos)
	link.Position = maxPos + 1

	query := `INSERT INTO links (user_id, title, url, icon, position, is_active) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, link.UserID, link.Title, link.URL, link.Icon, link.Position, link.IsActive)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	link.ID = id
	return nil
}

// GetByID retrieves a link by ID
func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*model.Link, error) {
	query := `SELECT id, user_id, title, url, icon, position, is_active, created_at FROM links WHERE id = ?`
	
	link := &model.Link{}
	var isActive int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID, &link.UserID, &link.Title, &link.URL, &link.Icon, &link.Position, &isActive, &link.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	link.IsActive = isActive == 1
	return link, nil
}

// GetByUserID retrieves all links for a user
func (r *LinkRepository) GetByUserID(ctx context.Context, userID int64) ([]model.Link, error) {
	query := `SELECT id, user_id, title, url, icon, position, is_active, created_at FROM links WHERE user_id = ? ORDER BY position ASC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		var isActive int
		if err := rows.Scan(&link.ID, &link.UserID, &link.Title, &link.URL, &link.Icon, &link.Position, &isActive, &link.CreatedAt); err != nil {
			return nil, err
		}
		link.IsActive = isActive == 1
		links = append(links, link)
	}
	return links, rows.Err()
}

// GetActiveByUserID retrieves active links (for public profile)
func (r *LinkRepository) GetActiveByUserID(ctx context.Context, userID int64) ([]model.Link, error) {
	query := `SELECT id, user_id, title, url, icon, position, is_active, created_at FROM links WHERE user_id = ? AND is_active = 1 ORDER BY position ASC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		var isActive int
		if err := rows.Scan(&link.ID, &link.UserID, &link.Title, &link.URL, &link.Icon, &link.Position, &isActive, &link.CreatedAt); err != nil {
			return nil, err
		}
		link.IsActive = isActive == 1
		links = append(links, link)
	}
	return links, rows.Err()
}

// Delete removes a link
func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM links WHERE id = ?", id)
	return err
}

// UpdatePositions for drag-reorder
func (r *LinkRepository) UpdatePositions(ctx context.Context, userID int64, positions map[int64]int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE links SET position = ? WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for linkID, position := range positions {
		if _, err := stmt.ExecContext(ctx, position, linkID, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
```

**🧠 Explain:**
> "is_active INTEGER - SQLite-ൽ boolean ഇല്ല. 0/1 use ചെയ്യുന്നു."
> "Transaction (tx) - Multiple updates ഒരുമിച്ച്. Fail ആയാൽ rollback."

---

## Create Analytics Repository

**⌨️ Create `internal/repository/analytics.go`:**
```go
package repository

import (
	"context"
	"database/sql"

	"linkbio/internal/model"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// RecordPageView records a profile page view
func (r *AnalyticsRepository) RecordPageView(ctx context.Context, userID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, event_type, referrer, user_agent) VALUES (?, 'page_view', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, referrer, userAgent)
	return err
}

// RecordLinkClick records a link click
func (r *AnalyticsRepository) RecordLinkClick(ctx context.Context, userID, linkID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, link_id, event_type, referrer, user_agent) VALUES (?, ?, 'link_click', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, linkID, referrer, userAgent)
	return err
}

// GetSummary returns analytics summary for dashboard
func (r *AnalyticsRepository) GetSummary(ctx context.Context, userID int64, days int) (*model.AnalyticsSummary, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN event_type = 'page_view' THEN 1 END) as views,
			COUNT(CASE WHEN event_type = 'link_click' THEN 1 END) as clicks
		FROM analytics 
		WHERE user_id = ? AND created_at >= datetime('now', '-' || ? || ' days')
	`
	
	summary := &model.AnalyticsSummary{}
	err := r.db.QueryRowContext(ctx, query, userID, days).Scan(&summary.TotalViews, &summary.TotalClicks)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
```

---

# PART 7: MIDDLEWARE (25:00 - 30:00)

## What is Middleware?

**📱 Read this:**
> "Middleware എന്താണ്? Request server-ൽ എത്തുമ്പോൾ handler-ന് മുമ്പ് run ആകുന്ന code. Security checkpoint പോലെ."

**🎯 Analogy:**
> "Airport പോലെ imagine ചെയ്യൂ. Plane-ൽ കയറുന്നതിന് മുമ്പ് security check, ticket check, boarding pass check. ഓരോന്നും ഒരു middleware."

## Create Middleware

**⌨️ Create `internal/middleware/middleware.go`:**
```go
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/time/rate"
)

type Middleware struct {
	log       *slog.Logger
	store     *sessions.CookieStore
	limiter   *rate.Limiter
	rateLimit int
}

func New(log *slog.Logger, sessionSecret string, rateLimit int) *Middleware {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return &Middleware{
		log:       log,
		store:     store,
		limiter:   rate.NewLimiter(rate.Limit(rateLimit), rateLimit*2),
		rateLimit: rateLimit,
	}
}

func (m *Middleware) Store() *sessions.CookieStore {
	return m.store
}

// Logger logs every request
func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		m.log.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start).String(),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Recovery catches panics
func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error("panic recovered", "error", err, "path", r.URL.Path)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(wrapped, r)
	})
}

// RateLimit prevents abuse
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

**🧠 Explain each middleware:**

> "Logger - ഓരോ request log ചെയ്യും. Debugging, monitoring."

> "Recovery - Code panic ആയാൽ (crash) catch ചെയ്ത് error response അയയ്ക്കും. Server crash ആകില്ല."

> "RateLimit - Too many requests block ചെയ്യും. DDoS protection."

---

## Create Auth Middleware

**⌨️ Add to `internal/middleware/auth.go`:**
```go
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
)

// Auth checks if user is logged in
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, "session")
		
		userID, ok := session.Values["user_id"].(int64)
		if !ok || userID == 0 {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		
		username, _ := session.Values["username"].(string)
		
		// Add user info to context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, usernameKey, username)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(ctx context.Context) int64 {
	if id, ok := ctx.Value(userIDKey).(int64); ok {
		return id
	}
	return 0
}

// UsernameFromContext extracts username from context
func UsernameFromContext(ctx context.Context) string {
	if name, ok := ctx.Value(usernameKey).(string); ok {
		return name
	}
	return ""
}
```

**🧠 Explain:**
> "Context - Request-ന്റെ കൂടെ data pass ചെയ്യാൻ. user_id, username middleware-ൽ set ചെയ്തു, handler-ൽ read ചെയ്യാം."

---

# PART 8: HANDLERS (30:00 - 40:00)

## Response Helper

**⌨️ Create `internal/pkg/response/response.go`:**
```go
package response

import (
	"html/template"
	"log/slog"
	"net/http"
)

type Responder struct {
	log       *slog.Logger
	templates *template.Template
}

func New(log *slog.Logger, templateDir string) (*Responder, error) {
	return &Responder{log: log}, nil
}

// Error sends an error response
func (r *Responder) Error(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

// HXRedirect sends HTMX redirect header
func (r *Responder) HXRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
}
```

---

## Template Helper

**⌨️ Create `internal/pkg/templates/templates.go`:**
```go
package templates

import (
	"html/template"
	"io"
	"strings"
)

// FuncMap returns template functions
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"multiply": func(a, b int) int { return a * b },
		"slice": func(s string, start, end int) string {
			if end > len(s) {
				end = len(s)
			}
			if start > len(s) {
				return ""
			}
			return s[start:end]
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
}

// Render parses and executes templates
func Render(w io.Writer, page string, data interface{}) error {
	tmpl, err := template.New("base.html").Funcs(FuncMap()).ParseFiles(
		"web/templates/layouts/base.html",
		"web/templates/pages/"+page,
	)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
```

---

## Handler Dependencies

**⌨️ Create `internal/handler/handler.go`:**
```go
package handler

import (
	"log/slog"

	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"

	"github.com/gorilla/sessions"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	Log           *slog.Logger
	Responder     *response.Responder
	Store         *sessions.CookieStore
	UserRepo      *repository.UserRepository
	LinkRepo      *repository.LinkRepository
	AnalyticsRepo *repository.AnalyticsRepository
}

// Handler groups all handlers
type Handler struct {
	Health    *HealthHandler
	Auth      *AuthHandler
	Link      *LinkHandler
	Dashboard *DashboardHandler
	Profile   *ProfileHandler
}

// New creates all handlers
func New(deps *Dependencies) *Handler {
	return &Handler{
		Health:    NewHealthHandler(deps),
		Auth:      NewAuthHandler(deps),
		Link:      NewLinkHandler(deps),
		Dashboard: NewDashboardHandler(deps),
		Profile:   NewProfileHandler(deps),
	}
}
```

---

## Health Handler

**⌨️ Create `internal/handler/health.go`:**
```go
package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler(deps *Dependencies) *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
```

---

## Auth Handler

**⌨️ Create `internal/handler/auth.go`:**
```go
package handler

import (
	"net/http"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/pkg/templates"
	"linkbio/internal/repository"

	"log/slog"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	log      *slog.Logger
	resp     *response.Responder
	store    *sessions.CookieStore
	userRepo *repository.UserRepository
}

func NewAuthHandler(deps *Dependencies) *AuthHandler {
	return &AuthHandler{
		log:      deps.Log,
		resp:     deps.Responder,
		store:    deps.Store,
		userRepo: deps.UserRepo,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "login.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.resp.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil {
		h.log.Error("database error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if user == nil {
		h.resp.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		h.resp.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	h.log.Info("user logged in", "user_id", user.ID)
	h.resp.HXRedirect(w, "/dashboard")
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render(w, "register.html", nil); err != nil {
		h.log.Error("template error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.resp.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || email == "" || password == "" {
		h.resp.Error(w, http.StatusBadRequest, "All fields are required")
		return
	}

	if len(password) < 6 {
		h.resp.Error(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	existing, _ := h.userRepo.GetByUsername(r.Context(), username)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Username already taken")
		return
	}

	existing, _ = h.userRepo.GetByEmail(r.Context(), email)
	if existing != nil {
		h.resp.Error(w, http.StatusConflict, "Email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  username,
		Theme:        "light",
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		h.log.Error("user creation error", "error", err)
		h.resp.Error(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	session, _ := h.store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(r, w)

	h.log.Info("user registered", "user_id", user.ID)
	h.resp.HXRedirect(w, "/dashboard")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
```

**🧠 Explain bcrypt:**
> "bcrypt - Password hashing algorithm. Plain password store ചെയ്യരുത്! Hashed version മാത്രം save ചെയ്യൂ."

---

**(Continue with remaining handlers, templates, and testing sections...)**

---

# PART 9: TESTING CHECKPOINT 1 (40:00)

## Test the Server

**📱 Read this:**
> "ഇപ്പോൾ വരെ എഴുതിയത് test ചെയ്യാം."

**⌨️ Create `cmd/server/main.go` and `internal/server/server.go` (as shown earlier)**

**⌨️ Run:**
```bash
go run ./cmd/server
```

**🎯 Expected output:**
```
INFO server starting port=8080
INFO database connected path=./data/linkbio.db
INFO database migrations completed
```

**🌐 Test in browser:**
- http://localhost:8080/health → `{"status":"ok"}`

> "Health check work ചെയ്യുന്നുണ്ടെങ്കിൽ, server ready ആണ്!"

---

# PART 10: TEMPLATES & FRONTEND (40:00 - 50:00)

## Base Layout

**📱 Read this:**
> "ഇനി HTML templates create ചെയ്യാം. Base layout - എല്ലാ pages-നും common structure."

**⌨️ Create `web/templates/layouts/base.html`:**
```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en" x-data="{ darkMode: localStorage.getItem('darkMode') === 'true' }" :class="{ 'dark': darkMode }">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}LinkBio{{end}}</title>
    
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            darkMode: 'class',
            theme: { extend: { fontFamily: { sans: ['Inter', 'system-ui', 'sans-serif'] } } }
        }
    </script>
    
    <link rel="stylesheet" href="/static/css/app.css">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
    <link href="https://unpkg.com/aos@2.3.1/dist/aos.css" rel="stylesheet">
    <script src="https://unpkg.com/aos@2.3.1/dist/aos.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/gsap/3.12.2/gsap.min.js"></script>
    
    {{block "head" .}}{{end}}
</head>
<body class="{{block "bodyClass" .}}bg-gray-50 dark:bg-gray-900{{end}} min-h-screen transition-colors duration-300">
    {{block "content" .}}{{end}}
    
    <script src="/static/js/app.js"></script>
    {{block "scripts" .}}{{end}}
</body>
</html>
{{end}}
```

**🧠 Explain:**
> "x-data - Alpine.js state. darkMode variable store ചെയ്യും."
> ":class - Dynamic class binding. darkMode true ആയാൽ 'dark' class add ചെയ്യും."
> "{{block}} - Override ചെയ്യാവുന്ന sections."

---

## Login Page — Alpine.js: Password Toggle & Loading State

**📱 Read this:**
> "Login page-ൽ Alpine.js use ചെയ്ത് രണ്ട് features add ചെയ്യാം. Password show/hide toggle-ഉം form submit loading state-ഉം."

**🎯 Why this matters:**
> "User password type ചെയ്യുമ്പോൾ typo ഉണ്ടോ എന്ന് check ചെയ്യാൻ show/hide toggle വേണം. Submit button-ൽ loading state ഉണ്ടെങ്കിൽ user-ന് feedback കിട്ടും - double click avoid ചെയ്യാം."

**⌨️ Key parts of `web/templates/pages/login.html`:**

```html
<!-- x-data: Alpine state — showPassword, loading -->
<div class="relative z-10 w-full max-w-md" x-data="{ showPassword: false, loading: false }">
```

> "x-data - Component-ന്റെ state define ചെയ്യുന്നു. showPassword false ആണ്, loading-ഉം false."

```html
<!-- Form: @submit sets loading, @htmx:after-request resets it -->
<form hx-post="/auth/login" 
      hx-target="#error-message"
      hx-swap="innerHTML"
      @submit="loading = true"
      @htmx:after-request.window="loading = false"
      class="space-y-5">
```

> "HTMX + Alpine.js together! @submit - form submit ചെയ്യുമ്പോൾ loading = true. HTMX response വരുമ്പോൾ @htmx:after-request loading = false ആക്കും."

```html
<!-- Password field with toggle -->
<div class="relative">
    <input :type="showPassword ? 'text' : 'password'" 
           id="password" name="password" required
           class="input-field w-full px-4 py-3.5 rounded-xl text-white placeholder-gray-500 pr-12"
           placeholder="••••••••">
    <button type="button" 
            @click="showPassword = !showPassword"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300">
        <svg x-show="!showPassword" class="w-5 h-5" ...><!-- Eye icon --></svg>
        <svg x-show="showPassword" x-cloak class="w-5 h-5" ...><!-- Eye-slash icon --></svg>
    </button>
</div>
```

**🧠 Explain:**
> ":type binding - showPassword true ആയാൽ type='text', false ആയാൽ type='password'. One line-ൽ toggle!"
> "@click - Button click ചെയ്യുമ്പോൾ showPassword flip ചെയ്യും."
> "x-show - Condition true ആയ icon മാത്രം show ചെയ്യും. x-cloak - Page load-ൽ flash ആകാതിരിക്കാൻ."

```html
<!-- Submit button with loading state -->
<button type="submit" 
        :disabled="loading"
        class="btn-primary w-full py-4 rounded-xl text-white font-semibold text-base mt-2"
        x-text="loading ? 'Signing in...' : 'Sign In'">
    Sign In
</button>
```

> ":disabled - loading true ആയാൽ button disabled ആകും. Double click prevent!"
> "x-text - Dynamic text. Loading ആയാൽ 'Signing in...', അല്ലെങ്കിൽ 'Sign In'."

**🎯 Alpine.js Concepts Used:**
| Directive | Purpose | Example |
|-----------|---------|---------|
| `x-data` | Component state | `{ showPassword: false }` |
| `:type` | Dynamic attribute | `showPassword ? 'text' : 'password'` |
| `@click` | Event handler | `showPassword = !showPassword` |
| `x-show` | Conditional display | Show/hide eye icon |
| `x-cloak` | Prevent flash | Hide until Alpine loads |
| `:disabled` | Dynamic disable | `loading` |
| `x-text` | Dynamic text | `loading ? '...' : '...'` |

---

## Register Page — Alpine.js: Live Preview, Password Strength & Loading

**📱 Read this:**
> "Register page-ൽ മൂന്ന് Alpine.js features ഉണ്ട്: Live username preview, password strength indicator, loading state."

**🎯 Analogy:**
> "Google account create ചെയ്യുമ്പോൾ password strength bar കണ്ടിട്ടുണ്ടല്ലോ? അത് തന്നെ നമ്മളും build ചെയ്യും — Alpine.js use ചെയ്ത്, JavaScript file ഇല്ലാതെ!"

**⌨️ Key parts of `web/templates/pages/register.html`:**

```html
<!-- Extended x-data: username, password, showPassword, loading -->
<div class="relative z-10 w-full max-w-md" 
     x-data="{ username: '', password: '', showPassword: false, loading: false }">
```

> "ഒരു x-data-ൽ multiple state variables. Alpine.js ഇത് automatically reactive ആക്കും."

**Feature 1: Live URL Preview (already exists)**
```html
<input type="text" id="username" name="username" x-model="username" ...>
<div x-show="username.length > 0" x-cloak>
    Your profile: linkbio.com/u/<span x-text="username"></span>
</div>
```

> "x-model - Two-way binding. Input-ൽ type ചെയ്യുമ്പോൾ username variable auto-update. Preview-ൽ real-time-ൽ URL കാണാം."

**Feature 2: Password Strength Indicator**
```html
<input :type="showPassword ? 'text' : 'password'" 
       x-model="password" ...>

<!-- Strength bar — only shows when user starts typing -->
<div x-show="password.length > 0" x-cloak class="mt-2">
    <div class="h-1.5 w-full bg-gray-700 rounded-full overflow-hidden">
        <div class="h-full rounded-full transition-all duration-300"
             :class="password.length < 6 ? 'bg-red-500' : password.length < 10 ? 'bg-yellow-500' : 'bg-green-500'"
             :style="'width: ' + Math.min(password.length * 10, 100) + '%'"></div>
    </div>
    <p class="text-xs mt-1"
       :class="password.length < 6 ? 'text-red-400' : password.length < 10 ? 'text-yellow-400' : 'text-green-400'"
       x-text="password.length < 6 ? 'Too short' : password.length < 10 ? 'Fair' : 'Strong'"></p>
</div>
```

**🧠 Explain step by step:**
> ":class - Ternary operator use ചെയ്ത് color change. 6-ൽ താഴെ red, 6-9 yellow, 10+ green."
> ":style - Bar width dynamic ആക്കും. ഓരോ character-നും 10% grow ചെയ്യും, maximum 100%."
> "x-text - Status text-ഉം dynamic. 'Too short', 'Fair', 'Strong'."

**🎯 Analogy:**
> "Traffic light പോലെ! Red = stop (too short), Yellow = caution (fair), Green = go (strong)."

**Feature 3: Loading State (same pattern as Login)**
```html
<form @submit="loading = true" @htmx:after-request.window="loading = false" ...>
    <button :disabled="loading"
            x-text="loading ? 'Creating account...' : 'Create Account'">
    </button>
</form>
```

---

## Dashboard — Alpine.js: Toggle Form, Transitions, Copy Link & HTMX Loading

**📱 Read this:**
> "Dashboard page-ൽ Alpine.js-ന്റെ power full-ൽ കാണാം. Form toggle, transitions, clipboard copy - എല്ലാം."

**⌨️ Key parts of `web/templates/pages/dashboard.html`:**

**Feature 1: Add Link Form Toggle with Transition**
```html
<main x-data="{ showAddForm: false }">
    <!-- Toggle button -->
    <button @click="showAddForm = !showAddForm">Add Link</button>
    
    <!-- Form with smooth transition -->
    <div x-show="showAddForm" x-cloak 
         x-transition:enter="transition ease-out duration-200"
         x-transition:enter-start="opacity-0 -translate-y-2"
         x-transition:enter-end="opacity-100 translate-y-0">
        <form hx-post="/api/v1/links" 
              hx-target="#links-list" 
              hx-swap="afterbegin"
              hx-indicator="find .htmx-indicator"
              @htmx:after-request="showAddForm = false; $el.reset()">
            ...
        </form>
    </div>
</main>
```

**🧠 Explain:**
> "x-transition - CSS animation Alpine.js handle ചെയ്യും. enter, enter-start, enter-end - 3 stages."
> "@htmx:after-request - HTMX response വന്നാൽ form close ചെയ്ത് reset ചെയ്യും. $el = current element (form)."

**🎯 Analogy:**
> "Drawer open/close ചെയ്യുന്നത് പോലെ. Button click → drawer slide down. Form submit → drawer slide up, form clear."

**Feature 2: Link Count Badge**
```html
<h2 class="text-lg font-semibold text-gray-900 dark:text-white">
    Your Links
    <span class="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400">
        {{len .Links}}
    </span>
</h2>
```

> "{{len .Links}} - Go template function. Server-side-ൽ link count calculate ചെയ്യും."

**Feature 3: HTMX Loading Spinner**
```html
<button type="submit" class="btn-primary ... inline-flex items-center gap-2">
    <svg class="w-4 h-4 animate-spin htmx-indicator" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
    </svg>
    Add Link
</button>
```

> "htmx-indicator class - HTMX request progress ആയാൽ show, complete ആയാൽ hide. CSS-ൽ opacity: 0 default, htmx-request class-ൽ opacity: 1."

**Feature 4: Copy Profile Link with Clipboard API**
```html
<button x-data="{ copied: false }"
        @click="navigator.clipboard.writeText(window.location.origin + '/u/{{.User.Username}}'); 
                copied = true; 
                setTimeout(() => copied = false, 2000)"
        :class="copied ? 'bg-green-500 text-white' : 'bg-gray-100 dark:bg-gray-800 ...'">
    <span x-text="copied ? '✓ Copied!' : 'Copy Profile Link'"></span>
</button>
```

**🧠 Explain:**
> "navigator.clipboard.writeText() - Browser Clipboard API. URL clipboard-ൽ copy ചെയ്യും."
> "setTimeout - 2 seconds-ന് ശേഷം copied = false ആക്കും. User-ന് visual feedback — green button-ൽ '✓ Copied!'"
> "Nested x-data - Parent-ന്റെ x-data-ൽ add ചെയ്യാതെ button-ന്റെ own state. Alpine.js scoping!"

---

## Profile Page — Alpine.js: Theme Toggle & Share Button

**📱 Read this:**
> "Public profile page-ൽ dark mode toggle-ഉം share button-ഉം add ചെയ്യാം."

**⌨️ Key parts of `web/templates/pages/profile.html`:**

**Feature 1: Independent Dark Mode (not saved to localStorage)**
```html
<div x-data="{ darkMode: false }" 
     :class="darkMode ? 'dark bg-gradient-to-b from-gray-900 to-gray-800' : 'bg-gradient-to-b from-gray-50 to-gray-100'">
```

> "Profile page-ന്റെ darkMode base layout-ന്റെ darkMode-ൽ നിന്ന് independent ആണ്. Visitor-ന് toggle ചെയ്യാം, save ചെയ്യില്ല."

**Feature 2: Theme Toggle with `x-if` Templates**
```html
<button @click="darkMode = !darkMode" :class="darkMode ? '...' : '...'">
    <template x-if="!darkMode">
        <span class="flex items-center gap-2">🌙 Dark Mode</span>
    </template>
    <template x-if="darkMode">
        <span class="flex items-center gap-2">☀️ Light Mode</span>
    </template>
</button>
```

> "x-if vs x-show: x-show CSS display:none use ചെയ്യും, element DOM-ൽ ഉണ്ടാകും. x-if element completely remove/add ചെയ്യും. Icon swap-ന് x-if better."

**Feature 3: Share Button with Clipboard Feedback**
```html
<button x-data="{ copied: false }"
        @click="navigator.clipboard.writeText(window.location.href); 
                copied = true; 
                setTimeout(() => copied = false, 2000)"
        :class="copied 
            ? 'bg-green-500 text-white' 
            : (darkMode ? 'bg-gray-800/80 text-gray-300' : 'bg-white/80 text-gray-600')">
    <template x-if="!copied">
        <span>📤 Share</span>
    </template>
    <template x-if="copied">
        <span>✓ Copied!</span>
    </template>
</button>
```

> "Parent-ന്റെ darkMode access ചെയ്യുന്നു, but own copied state-ഉം ഉണ്ട്. Alpine.js scope chain — inner component-ന് parent data access ഉണ്ട്!"

---

## 🎯 Alpine.js Quick Reference

| Directive | What it does | Used in |
|-----------|-------------|---------|
| `x-data` | Component state | All pages |
| `x-model` | Two-way input binding | Register (username, password) |
| `x-show` | Show/hide (CSS) | Login (eye icon), Dashboard (form) |
| `x-cloak` | Prevent load flash | With every x-show |
| `x-text` | Dynamic text content | Buttons, strength label |
| `x-if` | Add/remove from DOM | Profile (theme/share icons) |
| `:class` | Dynamic CSS classes | Strength bar colors, theme |
| `:type` | Dynamic input type | Password show/hide |
| `:disabled` | Dynamic disable | Submit buttons |
| `:style` | Dynamic inline style | Strength bar width |
| `@click` | Click handler | Toggles, clipboard |
| `@submit` | Form submit handler | Loading state |
| `x-transition` | Enter/leave animation | Dashboard add form |

> "ഈ 13 directives മാത്രം use ചെയ്ത് നമ്മൾ JavaScript file ഇല്ലാതെ 5 interactive features build ചെയ്തു. ഇതാണ് Alpine.js-ന്റെ power!"

---

# PART 11: HTMX MAGIC (50:00 - 55:00)

## What is HTMX?

**📱 Read this:**
> "HTMX - HTML over the wire. JavaScript എഴുതാതെ dynamic pages!"

**🎯 Analogy:**
> "Traditional way: API call → JSON → JavaScript parse → DOM update. HTMX way: API call → Server sends HTML → Browser swaps."

## HTMX Attributes

> "hx-post - Form submit ചെയ്യുമ്പോൾ POST request."
> "hx-target - Response എവിടെ insert ചെയ്യണം."
> "hx-swap - എങ്ങനെ insert ചെയ്യണം (innerHTML, outerHTML, afterbegin...)."

**Example:**
```html
<form hx-post="/api/v1/links" hx-target="#links-list" hx-swap="afterbegin">
```

> "ഈ form submit ചെയ്യുമ്പോൾ, server HTML return ചെയ്യും, HTMX അത് #links-list-ന്റെ തുടക്കത്തിൽ insert ചെയ്യും. Full page reload വേണ്ട!"

---

# PART 12: DEPLOY (55:00 - 58:00)

## Build for Production

**⌨️ Type:**
```bash
go build -o linkbio ./cmd/server
```

## Environment Variables

> "Production-ൽ .env file use ചെയ്യരുത്. System environment variables set ചെയ്യണം."

```bash
export ENV=production
export SESSION_SECRET=$(openssl rand -hex 32)
export DATABASE_PATH=/data/linkbio.db
```

---

# PART 13: RECAP (58:00 - 60:00)

## What We Built

> "ഇന്ന് നമ്മൾ build ചെയ്തത്:
> - Production-grade Go server with graceful shutdown
> - SQLite database with migrations
> - User authentication with sessions
> - HTMX for dynamic updates
> - Alpine.js for client-side state
> - Beautiful UI with Tailwind CSS"

## What's Next

> "Future videos-ൽ:
> - Google OAuth login
> - Stripe payments
> - Custom domains
> - Advanced analytics"

## Call to Action

> "Like, Subscribe, Share! Comment-ൽ നിങ്ങളുടെ questions ഇടൂ."

---

# 📝 QUICK REFERENCE COMMANDS

```bash
# Start development
go run ./cmd/server

# Build for production
go build -o linkbio ./cmd/server

# Run tests
go test ./...

# Check for issues
go vet ./...
```

---

# 🎯 TESTING CHECKPOINTS

| Time | What to Test | Expected Result |
|------|-------------|-----------------|
| 15:00 | `go build` | No errors |
| 25:00 | Server start | "database connected" log |
| 30:00 | /health | `{"status":"ok"}` |
| 40:00 | Register | Redirects to /dashboard |
| 45:00 | Add link | Link appears without reload |
| 50:00 | Public profile | Links displayed |
| 55:00 | Dark mode | Theme toggles |

---

# PART 14: TESTING & BENCHMARKS (Bonus Section)

## Why Testing Matters

**📱 Read this:**
> "Testing എന്തിന്? Production-ൽ bug കണ്ടെത്തുന്നതിനേക്കാൾ നല്ലത് development-ൽ കണ്ടെത്തുന്നത്. Testing നമ്മുടെ safety net ആണ്."

**🎯 Analogy:**
> "Car-ന്റെ brakes test ചെയ്യാതെ road-ൽ ഇറക്കുമോ? Code-ഉം അങ്ങനെ തന്നെ."

---

## Test Structure

**📱 Read this:**
> "Go-യിൽ test files `_test.go` suffix-ൽ end ആകണം. Same package-ൽ വെക്കാം."

```
internal/
├── handler/
│   ├── auth.go
│   ├── auth_test.go      ← Unit tests
│   ├── link.go
│   └── link_test.go
├── repository/
│   ├── user.go
│   └── user_test.go      ← Integration tests
```

---

## Step 1: Test Utilities

**⌨️ Create `internal/testutil/testutil.go`:**
```go
package testutil

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"linkbio/internal/repository"
)

// TestDB creates an in-memory database for testing
func TestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	if err := repository.Migrate(db, log); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// TestLogger creates a silent logger for testing
func TestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}
```

**🧠 Explain:**
> "t.Helper() - Test fail ആയാൽ correct line number കാണിക്കും."
> ":memory: - SQLite in-memory database. Test-ന് ശേഷം automatically delete ആകും."
> "t.Cleanup() - Test കഴിഞ്ഞാൽ auto-cleanup."

---

## Step 2: Repository Tests

**📱 Read this:**
> "ആദ്യം database operations test ചെയ്യാം. ഇത് integration tests ആണ് - real database use ചെയ്യുന്നു."

**⌨️ Create `internal/repository/user_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		DisplayName:  "Test User",
		Theme:        "light",
	}

	// Test: Create user
	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify: ID was set
	if user.ID == 0 {
		t.Error("Create() did not set user ID")
	}

	// Verify: User exists in database
	found, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if found == nil {
		t.Fatal("GetByID() returned nil")
	}
	if found.Username != user.Username {
		t.Errorf("Username = %v, want %v", found.Username, user.Username)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Setup: Create a user
	user := &model.User{
		Username:     "emailtest",
		Email:        "find@example.com",
		PasswordHash: "hash",
		DisplayName:  "Email Test",
		Theme:        "dark",
	}
	repo.Create(ctx, user)

	// Test cases
	tests := []struct {
		name      string
		email     string
		wantFound bool
	}{
		{
			name:      "existing email",
			email:     "find@example.com",
			wantFound: true,
		},
		{
			name:      "non-existing email",
			email:     "notfound@example.com",
			wantFound: false,
		},
		{
			name:      "empty email",
			email:     "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.GetByEmail(ctx, tt.email)
			if err != nil {
				t.Fatalf("GetByEmail() error = %v", err)
			}

			gotFound := found != nil
			if gotFound != tt.wantFound {
				t.Errorf("GetByEmail() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestUserRepository_DuplicateUsername(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Create first user
	user1 := &model.User{
		Username:     "duplicate",
		Email:        "user1@example.com",
		PasswordHash: "hash",
		DisplayName:  "User 1",
		Theme:        "light",
	}
	if err := repo.Create(ctx, user1); err != nil {
		t.Fatalf("Create() first user error = %v", err)
	}

	// Try to create second user with same username
	user2 := &model.User{
		Username:     "duplicate", // Same username!
		Email:        "user2@example.com",
		PasswordHash: "hash",
		DisplayName:  "User 2",
		Theme:        "light",
	}
	err := repo.Create(ctx, user2)

	// Should fail with unique constraint error
	if err == nil {
		t.Error("Create() should fail for duplicate username")
	}
}
```

**🧠 Explain Table-Driven Tests:**
> "tests := []struct{} - Table-driven testing. Multiple test cases ഒരു function-ൽ. Go community standard."
> "t.Run() - Sub-test. ഓരോ case independently run ആകും."

---

## Step 3: Link Repository Tests

**⌨️ Create `internal/repository/link_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func TestLinkRepository_Create(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup: Create a user first
	user := &model.User{
		Username:     "linktest",
		Email:        "link@test.com",
		PasswordHash: "hash",
		DisplayName:  "Link Test",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Test: Create link
	link := &model.Link{
		UserID:   user.ID,
		Title:    "My Website",
		URL:      "https://example.com",
		IsActive: true,
	}

	err := linkRepo.Create(ctx, link)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if link.ID == 0 {
		t.Error("Create() did not set link ID")
	}

	// Verify position was auto-set
	if link.Position != 1 {
		t.Errorf("Position = %v, want 1", link.Position)
	}
}

func TestLinkRepository_PositionAutoIncrement(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup
	user := &model.User{
		Username: "postest", Email: "pos@test.com",
		PasswordHash: "hash", DisplayName: "Pos Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	// Create 3 links
	for i := 1; i <= 3; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    "Link",
			URL:      "https://example.com",
			IsActive: true,
		}
		linkRepo.Create(ctx, link)

		if link.Position != i {
			t.Errorf("Link %d position = %v, want %v", i, link.Position, i)
		}
	}
}

func TestLinkRepository_GetActiveByUserID(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup
	user := &model.User{
		Username: "activetest", Email: "active@test.com",
		PasswordHash: "hash", DisplayName: "Active Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	// Create 2 active and 1 inactive link
	activeLink1 := &model.Link{UserID: user.ID, Title: "Active 1", URL: "https://a.com", IsActive: true}
	activeLink2 := &model.Link{UserID: user.ID, Title: "Active 2", URL: "https://b.com", IsActive: true}
	inactiveLink := &model.Link{UserID: user.ID, Title: "Inactive", URL: "https://c.com", IsActive: false}

	linkRepo.Create(ctx, activeLink1)
	linkRepo.Create(ctx, activeLink2)
	linkRepo.Create(ctx, inactiveLink)

	// Test: Get active links only
	activeLinks, err := linkRepo.GetActiveByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetActiveByUserID() error = %v", err)
	}

	if len(activeLinks) != 2 {
		t.Errorf("GetActiveByUserID() returned %d links, want 2", len(activeLinks))
	}

	// Verify all returned links are active
	for _, link := range activeLinks {
		if !link.IsActive {
			t.Errorf("GetActiveByUserID() returned inactive link: %s", link.Title)
		}
	}
}

func TestLinkRepository_Delete(t *testing.T) {
	db := testutil.TestDB(t)
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup
	user := &model.User{
		Username: "deltest", Email: "del@test.com",
		PasswordHash: "hash", DisplayName: "Del Test", Theme: "light",
	}
	userRepo.Create(ctx, user)

	link := &model.Link{UserID: user.ID, Title: "To Delete", URL: "https://del.com", IsActive: true}
	linkRepo.Create(ctx, link)

	// Test: Delete link
	err := linkRepo.Delete(ctx, link.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify: Link no longer exists
	found, _ := linkRepo.GetByID(ctx, link.ID)
	if found != nil {
		t.Error("Delete() did not remove link")
	}
}
```

---

## Step 4: Handler Tests

**📱 Read this:**
> "Handler tests - HTTP request simulate ചെയ്ത് response check ചെയ്യും. Unit testing ആണ്."

**⌨️ Create `internal/handler/health_test.go`:**
```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_Check(t *testing.T) {
	// Create handler
	h := &HealthHandler{}

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Call handler
	h.Check(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("status = %s, want ok", response["status"])
	}
}
```

**🧠 Explain:**
> "httptest.NewRequest() - Fake HTTP request create ചെയ്യുന്നു."
> "httptest.NewRecorder() - Response capture ചെയ്യുന്ന fake ResponseWriter."
> "Real server start ചെയ്യാതെ handler test ചെയ്യാം!"

---

## Step 5: Auth Handler Tests

**⌨️ Create `internal/handler/auth_test.go`:**
```go
package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/pkg/response"
	"linkbio/internal/repository"
	"linkbio/internal/testutil"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthHandler(t *testing.T) (*AuthHandler, *repository.UserRepository) {
	t.Helper()

	db := testutil.TestDB(t)
	log := testutil.TestLogger()

	userRepo := repository.NewUserRepository(db)
	resp, _ := response.New(log, "")
	store := sessions.NewCookieStore([]byte("test-secret-key-32-chars-minimum"))

	h := &AuthHandler{
		log:      log,
		resp:     resp,
		store:    store,
		userRepo: userRepo,
	}

	return h, userRepo
}

func TestAuthHandler_Login_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Setup: Create user with hashed password
	password := "testpass123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username:     "loginuser",
		Email:        "login@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Login User",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Create login request
	form := url.Values{}
	form.Add("email", "login@test.com")
	form.Add("password", password)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	// Call handler
	h.Login(rec, req)

	// Check: Should redirect via HX-Redirect header
	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Setup: Create user
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	user := &model.User{
		Username:     "wrongpass",
		Email:        "wrong@test.com",
		PasswordHash: string(hash),
		DisplayName:  "Wrong Pass",
		Theme:        "light",
	}
	userRepo.Create(ctx, user)

	// Try login with wrong password
	form := url.Values{}
	form.Add("email", "wrong@test.com")
	form.Add("password", "wrongpassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	// Check: Should return 401 Unauthorized
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("email", "notexists@test.com")
	form.Add("password", "anypassword")

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	form := url.Values{}
	form.Add("username", "newuser")
	form.Add("email", "new@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	// Check redirect
	if rec.Header().Get("HX-Redirect") != "/dashboard" {
		t.Errorf("HX-Redirect = %s, want /dashboard", rec.Header().Get("HX-Redirect"))
	}

	// Verify user was created
	user, _ := userRepo.GetByUsername(ctx, "newuser")
	if user == nil {
		t.Error("User was not created")
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	h, userRepo := setupAuthHandler(t)
	ctx := context.Background()

	// Create existing user
	existing := &model.User{
		Username:     "existing",
		Email:        "existing@test.com",
		PasswordHash: "hash",
		DisplayName:  "Existing",
		Theme:        "light",
	}
	userRepo.Create(ctx, existing)

	// Try to register with same username
	form := url.Values{}
	form.Add("username", "existing")
	form.Add("email", "different@test.com")
	form.Add("password", "password123")

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	h, _ := setupAuthHandler(t)

	form := url.Values{}
	form.Add("username", "shortpass")
	form.Add("email", "short@test.com")
	form.Add("password", "12345") // Less than 6 chars

	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	h.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
```

---

## Step 6: Benchmarks

**📱 Read this:**
> "Benchmarks - Code എത്ര fast ആണെന്ന് measure ചെയ്യാം. Performance optimization-ന്."

**⌨️ Create `internal/repository/benchmark_test.go`:**
```go
package repository

import (
	"context"
	"testing"

	"linkbio/internal/model"
	"linkbio/internal/testutil"
)

func BenchmarkUserRepository_Create(b *testing.B) {
	db := testutil.TestDB(&testing.T{})
	repo := NewUserRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &model.User{
			Username:     "benchuser" + string(rune(i)),
			Email:        "bench" + string(rune(i)) + "@test.com",
			PasswordHash: "hash",
			DisplayName:  "Bench User",
			Theme:        "light",
		}
		repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_GetByEmail(b *testing.B) {
	db := testutil.TestDB(&testing.T{})
	repo := NewUserRepository(db)
	ctx := context.Background()

	// Setup: Create user to find
	user := &model.User{
		Username:     "findme",
		Email:        "findme@test.com",
		PasswordHash: "hash",
		DisplayName:  "Find Me",
		Theme:        "light",
	}
	repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByEmail(ctx, "findme@test.com")
	}
}

func BenchmarkLinkRepository_GetByUserID(b *testing.B) {
	db := testutil.TestDB(&testing.T{})
	userRepo := NewUserRepository(db)
	linkRepo := NewLinkRepository(db)
	ctx := context.Background()

	// Setup: Create user with 10 links
	user := &model.User{
		Username: "linkbench", Email: "linkbench@test.com",
		PasswordHash: "hash", DisplayName: "Link Bench", Theme: "light",
	}
	userRepo.Create(ctx, user)

	for i := 0; i < 10; i++ {
		link := &model.Link{
			UserID:   user.ID,
			Title:    "Link",
			URL:      "https://example.com",
			IsActive: true,
		}
		linkRepo.Create(ctx, link)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		linkRepo.GetByUserID(ctx, user.ID)
	}
}
```

**🧠 Explain Benchmarks:**
> "b.N - Go automatically decide ചെയ്യും എത്ര times run ചെയ്യണമെന്ന്."
> "b.ResetTimer() - Setup time benchmark-ൽ include ചെയ്യരുത്."

---

## Step 7: HTTP Handler Benchmarks

**⌨️ Add to `internal/handler/health_test.go`:**
```go
func BenchmarkHealthHandler_Check(b *testing.B) {
	h := &HealthHandler{}

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		h.Check(rec, req)
	}
}
```

---

## Running Tests

**📱 Read this:**
> "Tests run ചെയ്യാൻ:"

**⌨️ Type:**
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/repository

# Run specific test function
go test -v -run TestUserRepository_Create ./internal/repository

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**🎯 Expected output:**
```
=== RUN   TestUserRepository_Create
--- PASS: TestUserRepository_Create (0.01s)
=== RUN   TestUserRepository_GetByEmail
=== RUN   TestUserRepository_GetByEmail/existing_email
=== RUN   TestUserRepository_GetByEmail/non-existing_email
=== RUN   TestUserRepository_GetByEmail/empty_email
--- PASS: TestUserRepository_GetByEmail (0.01s)
PASS
ok      linkbio/internal/repository     0.125s
```

---

## Running Benchmarks

**⌨️ Type:**
```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkUserRepository_GetByEmail ./internal/repository

# With memory allocation stats
go test -bench=. -benchmem ./...

# Run 5 times for accuracy
go test -bench=. -count=5 ./internal/repository
```

**🎯 Expected output:**
```
BenchmarkUserRepository_Create-8         	   10000	    120534 ns/op	    1024 B/op	      15 allocs/op
BenchmarkUserRepository_GetByEmail-8     	  100000	     15234 ns/op	     512 B/op	       8 allocs/op
BenchmarkLinkRepository_GetByUserID-8    	   50000	     28456 ns/op	    2048 B/op	      24 allocs/op
```

**🧠 Explain output:**
> "10000 - ഈ function 10000 times run ചെയ്തു."
> "120534 ns/op - ഓരോ operation-നും 120 microseconds."
> "1024 B/op - ഓരോ operation-നും 1KB memory allocate ചെയ്തു."
> "15 allocs/op - 15 memory allocations per operation."

---

## Test Coverage Goals

| Package | Target Coverage |
|---------|----------------|
| repository | 80%+ |
| handler | 70%+ |
| middleware | 60%+ |

**📱 Read this:**
> "100% coverage വേണ്ട. Critical paths cover ചെയ്താൽ മതി. Happy path + Error cases."

---

## Testing Best Practices

**📱 Read this:**
> "Testing tips:"

1. **AAA Pattern**: Arrange → Act → Assert
```go
// Arrange: Setup test data
user := &model.User{...}

// Act: Call the function
err := repo.Create(ctx, user)

// Assert: Check results
if err != nil { t.Fatal(err) }
```

2. **Table-Driven Tests**: Multiple cases, one function
3. **Test Isolation**: Each test independent
4. **Meaningful Names**: `TestUserRepository_Create_DuplicateEmail`
5. **Test Edge Cases**: Empty input, nil, boundaries

---

# 🎯 TESTING QUICK REFERENCE

```bash
# Quick test
go test ./...

# Verbose
go test -v ./...

# Coverage
go test -cover ./...

# Benchmark
go test -bench=. -benchmem ./...

# Race condition detection
go test -race ./...

# Specific test
go test -v -run TestName ./package

# Short mode (skip long tests)
go test -short ./...
```

---

**END OF SCRIPT**
