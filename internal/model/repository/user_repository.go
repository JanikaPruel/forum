package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"forum/internal/model"
	"forum/pkg/sqlite"
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
func (ur *UserRepository) GetUserByEmail(email string) *model.User {
	user := model.User{}
	err := ur.DB.SQLite.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE email=?", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error(err.Error())
		}
		return nil
	}

	return &user
}

// CreateUser
func (ur *UserRepository) CreateUser(user *model.User) (userID int, err error) {
	res, err := ur.DB.SQLite.Exec("INSERT INTO users (username, email, password) values(?, ?, ?)",
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
func (ur *UserRepository) CreateSession(userID int, expires time.Time) error {
	fmt.Println(userID, expires)
	res, err := ur.DB.SQLite.Exec("INSERT INTO sessions (user_id, expires) values(?, ?)",
		userID, expires)
	if err != nil {
		return err
	}

	fmt.Println(userID, expires)

	fmt.Println(res.RowsAffected())

	return nil
}

// RemoveSession
func (ur *UserRepository) RemoveSession(userID int) error {
	_, err := ur.DB.SQLite.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
