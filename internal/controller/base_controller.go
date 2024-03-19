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

type Date struct {
	Categories []model.Category
	Posts      []model.Post
	Comments   []model.Comment
	AuthUser   model.User
}

// MainController
func MainController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main")))

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}
