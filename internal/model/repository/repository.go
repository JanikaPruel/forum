package repository

import "forum/pkg/sqlite"

// Repository is a base repository
type Repository struct {
	URepo *UserRepository
	CRepo *CategoryRepository
}

// New Repository
func New(db *sqlite.Database) *Repository {
	return &Repository{
		URepo: NewUserRepository(db),
		CRepo: NewCateforyRepository(db),
	}
}
