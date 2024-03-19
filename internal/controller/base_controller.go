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
	Posts      []model.Post
	Comments   []model.Comment
	AuthUser   model.User
}

// MainController
func (ctl *BaseController) MainController(w http.ResponseWriter, r *http.Request) {
	// wd, err := os.Getwd() // /home/kooduser/Kood-tasks/forumV2 + ""
	// if err != nil {
	// 	slog.Error(err.Error())

	// }
	// fmt.Println("WORKDIR: ", wd)

	// auth -> login -> sign-up or sign-in
	// logout

	// categories
	categories, err := ctl.Repo.CRepo.GetAllCategories()
	if err != nil {
		slog.Error(err.Error())
	}
	// posts

	//

	//  <!-- Остальные пункты меню --> /home/kooduser/Kood-tasks/forumV2/internal/view/static/icons8-кот-64.png
	//   <!-- Остальные пункты меню --> /home/kooduser/Kood-tasks/forumV2/internal/controller/base_controller.gok

	// DATA
	data := Data{
		Categories: categories,
	}

	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main")))

	if err := tmp.Execute(w, data); err != nil {
		slog.Error(err.Error())
	}
}
