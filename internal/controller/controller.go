package controller

import (
	"forum/internal/model/repository"
	"forum/pkg/sqlite"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

const (
	mainPage  = "/internal/view/templates/main.html"
	loginPage = "/internal/view/templates/login.html"
	viewDir   = "/internal/view/"
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
	case "login.html", "login":
		tmplFilepath = wd + loginPage
	default:
		tmplFilepath = wd + viewDir
	}
	return tmplFilepath
}

// MainController
func MainController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main")))
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, nil); err != nil {
		slog.Error(err.Error())
	}
}
