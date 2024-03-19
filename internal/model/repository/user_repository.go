package repository

import (
	"fmt"
	"forum/internal/model"
	"forum/pkg/sqlite"
	"time"
)

// UserRepository
type UserRepository struct {
	DB *sqlite.Database
}

// NewUserRepository
func NewUserRepository(db *sqlite.Database) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// GetUserByEmail
func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	row := ur.DB.SQLite.QueryRow("SELECT id, username, email, password FROM users WHERE email=$1", email)

	user := &model.User{}
	if err := row.Scan(user.ID, user.Username, user.Email, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser
func (ur *UserRepository) CreateUser(user *model.User) (userID int, err error) {
	res, err := ur.DB.SQLite.Exec("INSERT INTO users (username, email, password) values($1, $2, $3)",
		user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Println(res.LastInsertId())
	fmt.Println(res.RowsAffected())

	return int(id), nil
}

// CreateSession
func (ur *UserRepository) CreateSession(sessionID string, userID int, expires time.Time) error {
	res, err := ur.DB.SQLite.Exec("INSERT INTO sessions (id, user_id, expires_at) values($1, $2, $3)",
		sessionID, userID, expires)
	if err != nil {
		return err
	}

	fmt.Println(sessionID, userID, expires)

	fmt.Println(res.RowsAffected())

	return nil
}
