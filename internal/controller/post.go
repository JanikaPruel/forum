package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"forum/internal/model"
)

// CreataPost | POST /posts
func (ctl *BaseController) CreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(GetWD() + "/web/templates/create_post.html"))

	user := ctl.GetAuthUser(r)
	logged := false
	if user != nil {
		logged = true
	}

	// categories
	categories, err := ctl.Repo.CRepo.GetAllCategories()
	if err != nil {
		slog.Error(err.Error())
	}

	err = r.ParseForm()
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	title := r.FormValue("post-title")
	content := r.FormValue("post-content")
	categoryIDsStr := r.Form["category-create-post[]"]
	if len(categoryIDsStr) == 0 {
		data := struct {
			Err        int
			Message    string
			Categories []*model.Category
			UserID     *model.User
			IsLoggedIn bool
		}{
			Err:        1,
			Message:    "Choose at least one category",
			Categories: categories,
			UserID:     user,
			IsLoggedIn: logged,
		}
		tmpl.Execute(w, data)
		return
	}

	categoryIDs := []int{}
	for _, categoryIDStr := range categoryIDsStr {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			data := struct {
				Err        int
				Message    string
				Categories []*model.Category
				UserID     *model.User
				IsLoggedIn bool
			}{
				Err:        1,
				Message:    "Choose at least one category",
				Categories: categories,
				UserID:     user,
				IsLoggedIn: logged,
			}
			tmpl.Execute(w, data)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	currentUser := ctl.GetAuthUser(r)

	post := model.Post{
		UserID:    currentUser.ID,
		Title:     title,
		Content:   content,
		Category:  categoryIDs,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now(),
	}
	if title == "" || content == "" {
		err := 4
		mes := "Title and content must be not empty, enter your data, and try again!"
		data := struct {
			Err     int
			Message string
		}{
			Err:     err,
			Message: mes,
		}
		tmpl.Execute(w, data)
		return
	}
	_, err = ctl.Repo.PRepo.CreatePost(post)
	if err != nil {
		slog.Error(err.Error())
		data := struct {
			Err        int
			Message    string
			Categories []*model.Category
			UserID     *model.User
			IsLoggedIn bool
		}{
			Err:        1,
			Message:    "Invalid data, try again!",
			Categories: categories,
			UserID:     user,
			IsLoggedIn: logged,
		}
		tmpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

// Function to search for posts by title and content
func (ctl *BaseController) SearchPosts(posts []*model.Post, searchQuery string, categoryFilters []int) []*model.Post {
	var filteredPosts []*model.Post

	for _, post := range posts {
		if len(categoryFilters) > 0 && !contains(post.Category, categoryFilters[0]) {
			continue
		}
		if strings.Contains(strings.ToLower(post.Title), strings.ToLower(searchQuery)) ||
			strings.Contains(strings.ToLower(post.Content), strings.ToLower(searchQuery)) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

// Function for generating a brief preview of the post content
func (ctl *BaseController) GeneratePreviewContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	return content[:maxChars] + "..."
}

func contains(slice []int, item int) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func (ctl *BaseController) SortPostsByDate(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
}
