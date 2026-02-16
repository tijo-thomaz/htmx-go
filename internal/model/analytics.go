package model

import "time"

// Analytics represents a tracking event
type Analytics struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LinkID    *int64    `json:"link_id,omitempty"` // nil for page views
	EventType string    `json:"event_type"`        // "page_view" or "link_click"
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// AnalyticsSummary holds aggregated analytics data
type AnalyticsSummary struct {
	TotalViews  int `json:"total_views"`
	TotalClicks int `json:"total_clicks"`
	LinkClicks  []LinkClickCount `json:"link_clicks"`
}

// LinkClickCount holds click count for a specific link
type LinkClickCount struct {
	LinkID int64  `json:"link_id"`
	Title  string `json:"title"`
	Clicks int    `json:"clicks"`
}
