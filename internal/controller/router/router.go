package router

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"forum/internal/controller"
)

// route structure <- serveMux - stdlib
type Router struct { // Структура для дефолтного или фреймоворкового мультиплексора. Со стороны архитектора такая структура задел на будущее, некоторый
	// фундамент и база который можно юзать. Это ООП подход и GREEN CASE.
	Mux *http.ServeMux
	// Chi Frame - work
	// Gin Frame - work
	// gorillaMux Frame - work
	Ctl *controller.BaseController
}

// New constractor - create a new router instance
func New(ctl *controller.BaseController) *Router {
	return &Router{
		Mux: http.NewServeMux(),
		Ctl: ctl,
	}
}

// InitRouter -
func (r *Router) InitRouter() {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}
	// internal/view/templates/main.html
	// /home/kooduser/Kood-tasks/forumV2/internal/view/templates/main.html
	fmt.Println("MESSAGE")
	fmt.Println(controller.GetTmplFilepath("main"))
	fmt.Println(controller.GetTmplFilepath("login"))
	fmt.Println("MESSAGE")
	// r.Mux.Handle("GET /templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(wd+"/internal/view/templates/"))))
	// r.Mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd+"/internal/view/static/"))))
	// r.Mux.Handle("GET /favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir(wd+"/internal/view/static/favicon/"))))

	r.Mux.Handle("GET /templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(wd+"/web/templates/"))))
	r.Mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd+"/web/static/"))))
	// r.Mux.Handle("GET /favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir(wd+"/internal/view/static/favicon/"))))

	r.Mux.HandleFunc("GET /login", controller.Login)
	r.Mux.HandleFunc("GET /sign-up", r.Ctl.SignUpPage)
	r.Mux.HandleFunc("POST /sign-up", r.Ctl.SignUp)
	r.Mux.HandleFunc("GET /sign-in", r.Ctl.SignInPage)
	r.Mux.HandleFunc("POST /sign-in", r.Ctl.SignIn)
	r.Mux.HandleFunc("GET /logout", r.Ctl.Logout)

	r.Mux.HandleFunc("GET /", r.Ctl.MainController)

	// categories
	r.Mux.HandleFunc("GET /categories", r.Ctl.MainController)
	r.Mux.HandleFunc("GET /admin/categories/{id}", r.Ctl.MainController)
	r.Mux.HandleFunc("POST /admin/categories", r.Ctl.CreatePost)
	r.Mux.HandleFunc("PUT /admin/categories/{id}", r.Ctl.MainController)
	r.Mux.HandleFunc("DELETE /admin/catefories/{id}", r.Ctl.MainController)

	// posts
	r.Mux.HandleFunc("GET /post", r.Ctl.ViewPostByID)
	r.Mux.HandleFunc("POST /posts", r.Ctl.CreatePost)
	// r.Mux.HandleFunc("PUT  /posts/id", r.Ctl.MainController)
	// r.Mux.HandleFunc("DELETE  /posts/{id}", r.Ctl.MainController)

	// comments
	r.Mux.HandleFunc("GET /comments", r.Ctl.MainController)
	r.Mux.HandleFunc("GET /comments/{id}", r.Ctl.MainController)
	r.Mux.HandleFunc("POST /comments", r.Ctl.MainController)
	r.Mux.HandleFunc("PUT /comments/{id}", r.Ctl.MainController)
	r.Mux.HandleFunc("DELETE /comments/{id}", r.Ctl.MainController)

	// likes
	r.Mux.HandleFunc("GET /comments/{id}/like", r.Ctl.MainController)
	r.Mux.HandleFunc("GET /post/{id}/like", r.Ctl.MainController)

	// dislikes
	r.Mux.HandleFunc("GET /comments/{id}/dislike", r.Ctl.MainController)
	r.Mux.HandleFunc("GET /post/{id}/dislike", r.Ctl.MainController)

} // 1.22

// endpoints:

// /sign-up, /sign-in

// /, /catefories

// /post/id, /add-comments, /remove-comment, /edit-comment /add-like, /remove-like, /add-dislike, /remove-dislike, /update-post

// assets - ./static, /view
