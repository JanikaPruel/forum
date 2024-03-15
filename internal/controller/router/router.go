package router

import (
	"forum/internal/controller"
	"net/http"
)

// route structure <- serveMux - stdlib
type Router struct { // Структура для дефолтного или фреймоворкового мультиплексора. Со стороны архитектора такая структура задел на будущее, некоторый
	// фундамент и база который можно юзать. Это ООП подход и GREEN CASE.
	Mux *http.ServeMux
	// Chi Frame - work
	// Gin Frame - work
	// gorillaMux Frame - work
}

// New constractor - create a new router instance
func New() *Router {
	return &Router{
		Mux: http.NewServeMux(),
	}
}

// InitRouter -
func (r *Router) InitRouter() {

	r.Mux.HandleFunc("POST /sign-up", controller.MainContrller)
	r.Mux.HandleFunc("POST /sign-in", controller.MainContrller)

	r.Mux.HandleFunc("GET /", controller.MainContrller)

	// categories
	r.Mux.HandleFunc("GET /categories", controller.MainContrller)
	r.Mux.HandleFunc("GET /admin/categories/{id}", controller.MainContrller)
	r.Mux.HandleFunc("POST /admin/categories", controller.MainContrller)
	r.Mux.HandleFunc("PUT /admin/categories/{id}", controller.MainContrller)
	r.Mux.HandleFunc("DELETE /admin/catefories/{id}", controller.MainContrller)

	// posts
	r.Mux.HandleFunc("GET  /posts", controller.MainContrller)
	r.Mux.HandleFunc("GET  /posts/{id}", controller.MainContrller)
	r.Mux.HandleFunc("POST /posts", controller.MainContrller)
	r.Mux.HandleFunc("PUT  /posts/id", controller.MainContrller)
	r.Mux.HandleFunc("DELETE  /posts/{id}", controller.MainContrller)

	// comments
	r.Mux.HandleFunc("GET /comments", controller.MainContrller)
	r.Mux.HandleFunc("GET /comments/{id}", controller.MainContrller)
	r.Mux.HandleFunc("POST /comments", controller.MainContrller)
	r.Mux.HandleFunc("PUT /comments/{id}", controller.MainContrller)
	r.Mux.HandleFunc("DELETE /comments/{id}", controller.MainContrller)

	// likes
	r.Mux.HandleFunc("GET /comments/{id}/like", controller.MainContrller)
	r.Mux.HandleFunc("GET /post/{id}/like", controller.MainContrller)

	// dislikes
	r.Mux.HandleFunc("GET /comments/{id}/dislike", controller.MainContrller)
	r.Mux.HandleFunc("GET /post/{id}/dislike", controller.MainContrller)

} // 1.22

// endpoints:

// /sign-up, /sign-in

// /, /catefories

// /post/id, /add-comments, /remove-comment, /edit-comment /add-like, /remove-like, /add-dislike, /remove-dislike, /update-post

// assets - ./static, /view
