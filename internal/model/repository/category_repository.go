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

// GetAllCategories
func (cr *CategoryRepository) GetAllCategories() ([]model.Category, error) {
	rows, err := cr.DB.SQLite.Query("SELECT id, name FROM categories")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()
	categories := []model.Category{}

	for rows.Next() {
		category := model.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			slog.Error(err.Error())
			continue
		}
		categories = append(categories, category)

	}

	return categories, nil
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

// GetCategoriesByPostID
func (cr *CategoryRepository) GetCategoriesByPostId(postID int) (categories []*model.Category, err error) {
	rows, err := cr.DB.SQLite.Query("SELECT categories.id, categories.name FROM post_categories INNER JOIN categories ON post_categories.category_id = categories.id WHERE post_categories.post_id = ?", postID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cat := model.Category{}
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			slog.Error(err.Error())
			continue
		}
		categories = append(categories, &cat)
	}
	return categories, nil
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
