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

// LinkCreateRequest is the input for creating a link
type LinkCreateRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Icon  string `json:"icon"`
}

// LinkUpdateRequest is the input for updating a link
type LinkUpdateRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Icon     string `json:"icon"`
	IsActive bool   `json:"is_active"`
}
