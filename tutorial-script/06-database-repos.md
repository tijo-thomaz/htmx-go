# Scene 6: Database & Repositories (23:00 - 28:00)

> 🎬 **Previous**: Models defined (Scene 5)
> 🎯 **Goal**: Database connection, migrations, all repository CRUD

---

## 🎥 Transition

> 📱 "Models ready. ഇനി actual database connection-ഉം queries-ഉം എഴുതാം."

---

## Database Connection & Migrations

> 📱 "Database connection code — connection pooling, WAL mode, migrations."

**⌨️ Create `internal/repository/db.go`:**
```go
package repository

import (
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

func NewDB(path string, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

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
		`CREATE INDEX IF NOT EXISTS idx_links_position ON links(user_id, position)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_user_id ON analytics(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_created_at ON analytics(created_at)`,
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

> 🧠 **Explain WAL mode:**
> 📱 "WAL = Write-Ahead Logging. Multiple readers ഒരേ time-ൽ read ചെയ്യാം, one writer write ചെയ്യാം. Default journal mode-നേക്കാൾ faster."

> 🧠 **Explain connection pooling:**
> 📱 "Connection pool - restaurant-ൽ plates wash and reuse ചെയ്യുന്നത് പോലെ. ഓരോ request-നും new connection create ചെയ്യാതെ pool-ൽ നിന്ന് reuse."

> 🧠 **Explain migrations:**
> 📱 "Migrations - app start ചെയ്യുമ്പോൾ tables auto-create ആകും. IF NOT EXISTS ഉള്ളതുകൊണ്ട് duplicate error വരില്ല."

---

## Migration SQL Files (Reference)

📱 **Narration**:
> "Database migrations code-ൽ inline ആണ്. But reference-ന് separate SQL files-ഉം create ചെയ്യാം."

⌨️ **Create `migrations/001_create_users.sql`:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    display_name TEXT DEFAULT '',
    bio TEXT DEFAULT '',
    avatar_url TEXT DEFAULT '',
    theme TEXT DEFAULT 'light',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

⌨️ **Create `migrations/002_create_links.sql`:**
```sql
CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    icon TEXT DEFAULT '',
    position INTEGER DEFAULT 0,
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

⌨️ **Create `migrations/003_create_analytics.sql`:**
```sql
CREATE TABLE IF NOT EXISTS analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    link_id INTEGER,
    event_type TEXT NOT NULL,
    referrer TEXT DEFAULT '',
    user_agent TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);
```

🧠 **Explain:**
> 📱 "ഈ files reference ആണ്. App-ൽ db.go inline migrations use ചെയ്യും. But SQL files keep ചെയ്യുന്നത് documentation-ന് നല്ലത്. Future-ൽ migration tool use ചെയ്യുമ്പോൾ ഈ files directly use ചെയ്യാം."

---

## User Repository

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

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (username, email, password_hash, display_name, bio, avatar_url, theme) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.DisplayName, user.Bio, user.AvatarURL, user.Theme,
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
	return user, err
}

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
	return user, err
}

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
	return user, err
}
```

> 🧠 **Explain:**
> 📱 "Repository pattern — database code ഒരു place-ൽ. Handler-ൽ SQL എഴുതരുത്. Separation of concerns."
> 📱 "Context — request cancel ആയാൽ query-ഉം cancel ആകും. Resources waste ചെയ്യില്ല."
> 📱 "sql.ErrNoRows — user not found ആയാൽ nil return ചെയ്യും, error അല്ല."

---

## Link Repository

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

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
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

func (r *LinkRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM links WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

func (r *LinkRepository) Update(ctx context.Context, link *model.Link) error {
	_, err := r.db.ExecContext(ctx, `UPDATE links SET title = ?, url = ?, icon = ?, is_active = ? WHERE id = ?`,
		link.Title, link.URL, link.Icon, link.IsActive, link.ID)
	return err
}

func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM links WHERE id = ?", id)
	return err
}

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

> 🧠 **Explain:**
> 📱 "COALESCE(MAX(position), 0) — auto-position. New link always goes to the end."
> 📱 "is_active INTEGER — SQLite-ൽ boolean ഇല്ല. 0 = false, 1 = true."
> 📱 "Transaction (tx) — reorder-ൽ multiple updates ഒരുമിച്ച്. ഒന്ന് fail ആയാൽ rollback — എല്ലാം undo."

> 🎯 **Analogy:**
> 📱 "Transaction = bank transfer. Account A-ൽ നിന്ന് B-ലേക്ക് transfer. A-ൽ debit ചെയ്ത ശേഷം B-ൽ credit fail ആയാൽ? Rollback — A-ൽ money back."

---

## Analytics Repository

**⌨️ Create `internal/repository/analytics.go`:**
```go
package repository

import (
	"context"
	"database/sql"
	"time"

	"linkbio/internal/model"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) RecordPageView(ctx context.Context, userID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, event_type, referrer, user_agent) VALUES (?, 'page_view', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, referrer, userAgent)
	return err
}

func (r *AnalyticsRepository) RecordLinkClick(ctx context.Context, userID, linkID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, link_id, event_type, referrer, user_agent) VALUES (?, ?, 'link_click', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, linkID, referrer, userAgent)
	return err
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context, userID int64, days int) (*model.AnalyticsSummary, error) {
	since := time.Now().AddDate(0, 0, -days)
	summary := &model.AnalyticsSummary{}

	// Query 1: Total page views
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM analytics WHERE user_id = ? AND event_type = 'page_view' AND created_at >= ?`,
		userID, since,
	).Scan(&summary.TotalViews)
	if err != nil {
		return nil, err
	}

	// Query 2: Total link clicks
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM analytics WHERE user_id = ? AND event_type = 'link_click' AND created_at >= ?`,
		userID, since,
	).Scan(&summary.TotalClicks)
	if err != nil {
		return nil, err
	}

	// Query 3: Per-link click breakdown
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.link_id, l.title, COUNT(*) as clicks
		FROM analytics a
		JOIN links l ON a.link_id = l.id
		WHERE a.user_id = ? AND a.event_type = 'link_click' AND a.created_at >= ?
		GROUP BY a.link_id
		ORDER BY clicks DESC
	`, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lc model.LinkClickCount
		if err := rows.Scan(&lc.LinkID, &lc.Title, &lc.Clicks); err != nil {
			return nil, err
		}
		summary.LinkClicks = append(summary.LinkClicks, lc)
	}

	return summary, rows.Err()
}
```

> 🧠 **Explain GetSummary — 3 queries:**
> 📱 "Query 1 — total page views. Profile page visit ചെയ്യുമ്പോൾ record ആകുന്നത്."
> 📱 "Query 2 — total link clicks. Link click ചെയ്യുമ്പോൾ record ആകുന്നത്."
> 📱 "Query 3 — per-link breakdown. JOIN ചെയ്ത് link title-ഉം edukkunnu. GROUP BY link_id — ഓരോ link-നും count. ORDER BY clicks DESC — most popular first."

> 📱 "time.Now().AddDate(0, 0, -days) — last N days. Dashboard-ൽ 28 days pass ചെയ്യും."

---

> 🎥 **Transition:** "Database layer complete! ഇനി middleware — security guard."
