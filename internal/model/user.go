package model

import "time"

// User represents a registered user
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	Theme        string    `json:"theme"`
	CreatedAt    time.Time `json:"created_at"`
}
