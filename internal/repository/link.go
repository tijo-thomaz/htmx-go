package repository

import (
	"context"
	"database/sql"

	"linkbio/internal/model"
)

// LinkRepository handles link database operations
type LinkRepository struct {
	db *sql.DB
}

// NewLinkRepository creates a new LinkRepository
func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

// Create inserts a new link
func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	// Get next position
	var maxPos int
	r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(position), 0) FROM links WHERE user_id = ?", link.UserID).Scan(&maxPos)
	link.Position = maxPos + 1

	query := `
		INSERT INTO links (user_id, title, url, icon, position, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		link.UserID,
		link.Title,
		link.URL,
		link.Icon,
		link.Position,
		link.IsActive,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	link.ID = id

	return nil
}

// GetByID retrieves a link by ID
func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*model.Link, error) {
	query := `
		SELECT id, user_id, title, url, icon, position, is_active, created_at
		FROM links WHERE id = ?
	`
	link := &model.Link{}
	var isActive int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&link.ID,
		&link.UserID,
		&link.Title,
		&link.URL,
		&link.Icon,
		&link.Position,
		&isActive,
		&link.CreatedAt,
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

// GetByUserID retrieves all links for a user ordered by position
func (r *LinkRepository) GetByUserID(ctx context.Context, userID int64) ([]model.Link, error) {
	query := `
		SELECT id, user_id, title, url, icon, position, is_active, created_at
		FROM links WHERE user_id = ?
		ORDER BY position ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		var isActive int
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Title,
			&link.URL,
			&link.Icon,
			&link.Position,
			&isActive,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}
		link.IsActive = isActive == 1
		links = append(links, link)
	}

	return links, rows.Err()
}

// GetActiveByUserID retrieves active links for a user (for public profile)
func (r *LinkRepository) GetActiveByUserID(ctx context.Context, userID int64) ([]model.Link, error) {
	query := `
		SELECT id, user_id, title, url, icon, position, is_active, created_at
		FROM links WHERE user_id = ? AND is_active = 1
		ORDER BY position ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		var isActive int // SQLite stores bool as int
		if err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Title,
			&link.URL,
			&link.Icon,
			&link.Position,
			&isActive,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}
		link.IsActive = isActive == 1
		links = append(links, link)
	}

	return links, rows.Err()
}

// Update updates a link
func (r *LinkRepository) Update(ctx context.Context, link *model.Link) error {
	query := `
		UPDATE links 
		SET title = ?, url = ?, icon = ?, is_active = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		link.Title,
		link.URL,
		link.Icon,
		link.IsActive,
		link.ID,
	)
	return err
}

// Delete removes a link
func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM links WHERE id = ?", id)
	return err
}

// UpdatePositions updates link positions (for drag-reorder)
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

func (r *LinkRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM links WHERE user_id = ?", userID).Scan(&count)
	return count, err
}
