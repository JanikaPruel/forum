package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

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

// New
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

type Data struct {
	Categories []model.Category
	Posts      []*model.Post
	Comments   []model.Comment
	UserID     *model.User
	Cookie     *http.Cookie
}

// MainController
func (ctl *BaseController) MainController(w http.ResponseWriter, r *http.Request) {
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

	cookie, _ := r.Cookie("sissionID")

	user := ctl.GetAuthUser(r)
	if user == nil {
		slog.Error("user is nil, underfind")
	}
	// us := &model.User{
	// 		ID:       1,
	// 		Username: "USER",
	// 		Email:    "user@user.com",
	// 		Password: "asdajslkdjqlkwejlqkwje",
	// }

	// DATA
	data := Data{
		Categories: categories,
		Posts:      posts,
		UserID:     user,
		Cookie:     cookie,
	}

	wd, _ := os.Getwd()

	tmp := template.Must(template.ParseFiles(wd + "/web/templates/main.html"))

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
	}
}
