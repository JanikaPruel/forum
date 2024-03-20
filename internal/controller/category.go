package controller

import (
	"log/slog"

	"forum/internal/model"
)

// GetCategories
func (ctl *BaseController) GetCategories() (categories *[]model.Category, err error) {
	// get all category from db

	// categories, err = ctl.Repo.CRepo.GetCategories()
	if err != nil {
		slog.Error(err.Error())
	}

	return nil, nil
}

// GetCategory by ID
func (ctl *BaseController) GetCategoryById(id int) *model.Category {
	categories, _ := ctl.Repo.CRepo.GetAllCategories()
	for _, category := range categories {
		if category.ID == id {
			return category
		}
	}
	return nil
}

// CreateCategory

// Update Category

// RemoveCategory

// FilterByCategory posts by category
func FilterByCategory(posts []*model.Post, categoryIDs []int) []*model.Post {
	filteredPosts := make([]*model.Post, 0)
	for _, post := range posts {
		for _, categoryID := range categoryIDs {
			if contains(post.Category, categoryID) {
				filteredPosts = append(filteredPosts, post)
				break
			}
		}
	}
	return filteredPosts
}
