package repository

import "forum/pkg/sqlite"

// CommentRepository
type CommentRepository struct {
	DB *sqlite.Database
}

// NewCommentRepository
func NewCommentRepository(db *sqlite.Database) *CommentRepository {
	return &CommentRepository{
		DB: db,
	}
}