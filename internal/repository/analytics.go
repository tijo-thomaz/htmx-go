package repository

import (
	"context"
	"database/sql"
	"time"

	"linkbio/internal/model"
)

// AnalyticsRepository handles analytics database operations
type AnalyticsRepository struct {
	db *sql.DB
}

// NewAnalyticsRepository creates a new AnalyticsRepository
func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// RecordPageView records a page view event
func (r *AnalyticsRepository) RecordPageView(ctx context.Context, userID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, event_type, referrer, user_agent) VALUES (?, 'page_view', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, referrer, userAgent)
	return err
}

// RecordLinkClick records a link click event
func (r *AnalyticsRepository) RecordLinkClick(ctx context.Context, userID, linkID int64, referrer, userAgent string) error {
	query := `INSERT INTO analytics (user_id, link_id, event_type, referrer, user_agent) VALUES (?, ?, 'link_click', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, linkID, referrer, userAgent)
	return err
}

// GetSummary retrieves analytics summary for a user within the given days
func (r *AnalyticsRepository) GetSummary(ctx context.Context, userID int64, days int) (*model.AnalyticsSummary, error) {
	since := time.Now().AddDate(0, 0, -days)

	summary := &model.AnalyticsSummary{}

	// Get total views
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM analytics WHERE user_id = ? AND event_type = 'page_view' AND created_at >= ?`,
		userID, since,
	).Scan(&summary.TotalViews)
	if err != nil {
		return nil, err
	}

	// Get total clicks
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM analytics WHERE user_id = ? AND event_type = 'link_click' AND created_at >= ?`,
		userID, since,
	).Scan(&summary.TotalClicks)
	if err != nil {
		return nil, err
	}

	// Get clicks per link
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
