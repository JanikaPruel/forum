package controller

import (
	"log/slog"

	"forum/internal/model"
)

// GetCategories
func (ctl *BaseController) GetCategories() (categories *[]model.Category, err error) {
	// get all category from db
	
	categories, err = ctl.Repo.CRepo.GetCategories()
	if err != nil {
		slog.Error(err.Error())
	}

	return nil, nil
}

// GetCategory by ID

// CreateCategory

// Update Category

// RemoveCategory
