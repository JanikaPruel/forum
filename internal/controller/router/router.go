package router

import "net/http"

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

	r.Mux.HandleFunc("POST /sign-up")
	r.Mux.HandleFunc("POST /sign-in")

	r.Mux.HandleFunc("GET /")

	// categories
	r.Mux.HandleFunc("GET /categories")
	r.Mux.HandleFunc("GET /admin/categories/{id}")
	r.Mux.HandleFunc("POST /admin/categories")
	r.Mux.HandleFunc("PUT /admin/categories/{id}")
	r.Mux.HandleFunc("DELETE /admin/catefories/{id}")

	// posts
	r.Mux.HandleFunc("GET  /posts")
	r.Mux.HandleFunc("GET  /posts/{id}")
	r.Mux.HandleFunc("POST /posts")
	r.Mux.HandleFunc("PUT  /posts/id")
	r.Mux.HandleFunc("DELETE  /posts/{id}")

	// comments
	r.Mux.HandleFunc("GET /comments")
	r.Mux.HandleFunc("GET /comments/{id}")
	r.Mux.HandleFunc("POST /comments")
	r.Mux.HandleFunc("PUT /comments/{id}")
	r.Mux.HandleFunc("DELETE /comments/{id}")

	// likes
	r.Mux.HandleFunc("GET /comments/{id}/like")
	r.Mux.HandleFunc("GET /post/{id}/like")

	// dislikes
	r.Mux.HandleFunc("GET /comments/{id}/dislike")
	r.Mux.HandleFunc("GET /post/{id}/dislike")

} // 1.22

// endpoints:

// /sign-up, /sign-in

// /, /catefories

// /post/id, /add-comments, /remove-comment, /edit-comment /add-like, /remove-like, /add-dislike, /remove-dislike, /update-post

// assets - ./static, /view
