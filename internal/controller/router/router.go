package router

import (
	"fmt"
	"forum/internal/controller"
	"log/slog"
	"net/http"
	"os"
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
	r.Mux.Handle("GET /templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir(wd+"/internal/view/templates/"))))
	r.Mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd+"/internal/view/static/"))))
	r.Mux.Handle("GET /favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir(wd+"/internal/view/static/favicon/"))))

	r.Mux.HandleFunc("GET /login", controller.Login)
	r.Mux.HandleFunc("POST /sign-up", r.Ctl.SignUp)
	r.Mux.HandleFunc("POST /sign-in", controller.SignIn)
	r.Mux.HandleFunc("GET /logout", controller.Logout)

	r.Mux.HandleFunc("GET /", controller.MainController)

	// categories
	r.Mux.HandleFunc("GET /categories", controller.MainController)
	r.Mux.HandleFunc("GET /admin/categories/{id}", controller.MainController)
	r.Mux.HandleFunc("POST /admin/categories", controller.MainController)
	r.Mux.HandleFunc("PUT /admin/categories/{id}", controller.MainController)
	r.Mux.HandleFunc("DELETE /admin/catefories/{id}", controller.MainController)

	// posts
	r.Mux.HandleFunc("GET  /posts", controller.MainController)
	r.Mux.HandleFunc("GET  /posts/{id}", controller.MainController)
	r.Mux.HandleFunc("POST /posts", controller.MainController)
	r.Mux.HandleFunc("PUT  /posts/id", controller.MainController)
	r.Mux.HandleFunc("DELETE  /posts/{id}", controller.MainController)

	// comments
	r.Mux.HandleFunc("GET /comments", controller.MainController)
	r.Mux.HandleFunc("GET /comments/{id}", controller.MainController)
	r.Mux.HandleFunc("POST /comments", controller.MainController)
	r.Mux.HandleFunc("PUT /comments/{id}", controller.MainController)
	r.Mux.HandleFunc("DELETE /comments/{id}", controller.MainController)

	// likes
	r.Mux.HandleFunc("GET /comments/{id}/like", controller.MainController)
	r.Mux.HandleFunc("GET /post/{id}/like", controller.MainController)

	// dislikes
	r.Mux.HandleFunc("GET /comments/{id}/dislike", controller.MainController)
	r.Mux.HandleFunc("GET /post/{id}/dislike", controller.MainController)

} // 1.22

// endpoints:

// /sign-up, /sign-in

// /, /catefories

// /post/id, /add-comments, /remove-comment, /edit-comment /add-like, /remove-like, /add-dislike, /remove-dislike, /update-post

// assets - ./static, /view
