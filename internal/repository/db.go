package repository

import (
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

// NewDB creates a new database connection with connection pooling
func NewDB(path string, log *slog.Logger) (*sql.DB, error) {
	// Open database with WAL mode and busy timeout for concurrent access
	db, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	// Connection pool settings for concurrent users
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Info("database connected", "path", path)

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *sql.DB, log *slog.Logger) error {
	migrations := []string{
		// Users table
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

		// Links table
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

		// Analytics table
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

		// Indexes for performance
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
