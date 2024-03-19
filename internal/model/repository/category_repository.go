package repository

import (
	"log/slog"

	"forum/internal/model"
	"forum/pkg/sqlite"
)

// Category repository
type CategoryRepository struct {
	DB *sqlite.Database
}

// New
func NewCateforyRepository(db *sqlite.Database) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

// GetCategories
func (cr *CategoryRepository) GetAllCategories() (*[]model.Category, error) {
	rows, err := cr.DB.SQLite.Query("SELECT * FROM categories")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	categories := []model.Category{}

	if rows.Next() {
		category := model.Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		categories = append(categories, category)

	}

	return &categories, nil
}

// GetCategory by ID
func (cr *CategoryRepository) GetCategoryByID(categoryID int) (category model.Category, err error) {
	err = cr.DB.SQLite.QueryRow("SELECT * FROM categories WHERE id=?", categoryID).
		Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return category, err
	}

	return category, nil
}

// CreateCategory
func (cr *CategoryRepository) CreateCategory(category *model.Category) error {
	_, err := cr.DB.SQLite.Exec("INSERT INTO categories(name) values(?)", category.Name)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	// slog.Info("Category: ", "ID: ",res.LastInsertId())
	return nil
}

// Update Category by ID
func (cr *CategoryRepository) UpdateCategory(categoryID int, categoryName string) error {
	_, err := cr.DB.SQLite.Exec("UPDATE categories SET name=? WHERE id=?", categoryName, categoryID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// DeleteCategory by ID
func (cr *CategoryRepository) DeleteCategory(categoryID int) error {
	_, err := cr.DB.SQLite.Exec("DELETE FROM categories WHERE id=?", categoryID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
