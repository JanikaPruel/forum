package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"forum/internal/model"
	"forum/internal/model/repository"
	"forum/pkg/sqlite"
)

const (
	mainPage       = "/internal/view/templates/main.html"
	loginPage      = "/internal/view/templates/login.html"
	CategoriesPage = "/internal/view/templates/categories.html"
	viewDir        = "/internal/view/"
)

type BaseController struct {
	Repo *repository.Repository
}

func New(db *sqlite.Database) *BaseController {
	return &BaseController{
		Repo: repository.New(db),
	}
}

func GetTmplFilepath(tmplName string) (tmplFilepath string) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	switch tmplName {
	case "main.html", "main":
		tmplFilepath = wd + mainPage
	case "categories.html", "categories", "category":
		tmplFilepath = wd + CategoriesPage
	case "login.html", "login":
		tmplFilepath = wd + loginPage
	default:
		tmplFilepath = wd + viewDir
	}
	return tmplFilepath
}

type Reaction struct {
	ID     int    `json:"id"`
	MarkID string `json:"mark_id"`
	UserID int    `json:"user_id"`
	Mark   string `json:"mark"`
}

type Data struct {
	Categories []*model.Category
	// Post       []*model.Post
	Posts    []*PostWithComments
	Comments []model.Comment
	UserID   *model.User
	Cookie   *http.Cookie
	Logged   bool
}

type PostWithComments struct {
	Post              *model.Post
	Author            *model.User
	CreatedAt         time.Time
	Comments          []*model.Comment
	PostLikes         int
	PostDislikes      int
	PreviewContent    string
	PreviewCategories string
	AuthUserID        int
	Categories        []*model.Category
}

// MainController
func (ctl *BaseController) MainController(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	tmp := template.Must(template.ParseFiles(wd + "/web/templates/main.html"))

	// Error test
	// if r.URL.Path != "/" {
	// 	ctl.ErrorController(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	// 	return
	// }

	// auth -> login -> sign-up or sign-in
	// logout

	// categories
	categories, err := ctl.Repo.CRepo.GetAllCategories()
	if err != nil {
		slog.Error(err.Error())
	}
	// posts
	posts, err := ctl.Repo.PRepo.GetAllPosts()
	if err != nil {
		slog.Error(err.Error())
	}

	user := ctl.GetAuthUser(r)
	logged := true
	if user == nil {
		logged = false
	}
	// us := &model.User{
	// 		ID:       1,
	// 		Username: "USER",
	// 		Email:    "user@user.com",
	// 		Password: "asdajslkdjqlkwejlqkwje",
	// }

	// revieve search and filters params
	categoriesID := r.URL.Query()["categories"]
	searchQuery := r.URL.Query().Get("search")

	categoryFilters := []int{}

	for _, categoryID := range categoriesID {
		categoryFilter, err := strconv.Atoi(categoryID)
		if err != nil {
			slog.Error(err.Error())
			ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		// add
		categoryFilters = append(categoryFilters, categoryFilter)
	}

	if len(categoryFilters) > 0 {
		posts = FilterByCategory(posts, categoryFilters)
	}

	if len(categoryFilters) > 0 || searchQuery != "" {
		posts = ctl.SearchPosts(posts, searchQuery, categoryFilters)
	}

	ctl.SortPostsByDate(posts)

	sortOption := r.URL.Query().Get("sort")
	if sortOption == "likes" {
		ctl.FilterByLikes(posts)
	} else if sortOption == "dislikes" {
		ctl.FilterByDislikes(posts)
	}

	likedPosts := []*model.TotalLikesPost{}
	dislikedPosts := []*model.TotalDislikesPost{}

	postsWithComments := make([]*PostWithComments, len(posts))

	for i, post := range posts {
		previewContent := ctl.GeneratePreviewContent(post.Content, 100)
		categories, _ := ctl.Repo.CRepo.GetCategoriesByPostID(post.ID)
		comments, _ := ctl.Repo.ComRepo.GetCommentsByPostID(post.ID)
		postLikes := len(ctl.Repo.PRepo.GetAllPostLikesByUserID(post.UserID))
		postDislikes := len(ctl.Repo.PRepo.GetAllPostDislikesByUserIDs(post.UserID))
		author, _ := ctl.Repo.URepo.GetUserByID(post.UserID)
		authUser := ctl.GetAuthUser(r)
		authUserID := 0
		if authUser != nil {
			authUserID = authUser.ID
		}
		isLiked := false
		for _, likePost := range likedPosts {
			if likePost.PostID == post.ID {
				isLiked = true
			}
		}
		post.IsLikedByAuthUser = isLiked
		var categoryNames []string
		for _, categoryId := range post.Category {
			category := ctl.GetCategoryById(categoryId)
			categoryNames = append(categoryNames, category.Name)
		}

		categoriesString := strings.Join(categoryNames, ", ")
		isDisliked := false
		for _, dislikePost := range dislikedPosts {
			if dislikePost.PostID == post.ID {
				isDisliked = true
			}
		}
		previewCategories := generatePreviewContent(categoriesString, 50)
		post.IsDislikedByAuthUser = isDisliked
		postsWithComments[i] = &PostWithComments{
			Post:              post,
			Author:            author,
			CreatedAt:         post.CreatedAt,
			Comments:          comments,
			PostLikes:         postLikes,
			PostDislikes:      postDislikes,
			PreviewContent:    previewContent,
			PreviewCategories: previewCategories,
			AuthUserID:        authUserID,
			Categories:        categories,
		}
	}
	resMap := make(template.FuncMap)
	resMap["Split"] = strings.Split

	data := Data{
		UserID:     user,
		Posts:      postsWithComments,
		Categories: categories,
		Logged:     logged,
	}

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
	}
}

func generatePreviewContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	return content[:maxChars] + "..."
}
