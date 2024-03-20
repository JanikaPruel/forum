package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"forum/internal/model"
)

// CreateComment handles the HTTP POST request for creating a comment.
//
// It expects the request to have the "post_id" and "comment_content" fields in
// its form data. If the request is not a POST request, or if the current user
// is not authenticated, it returns an error. It also checks if the provided
// post_id is a valid integer, and if the comment_content is not empty.
//
// The function retrieves the list of all posts and finds the corresponding post
// for the provided post_id. It then creates a new comment with the current user's
// information, the post's username, the provided post_id, and the comment_content.
// Finally, it redirects the user to the post's page.
func (ctl *BaseController) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		// Return an error if the request method is not POST
		ctl.ErrorController(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Get the authenticated user
	currentUser := ctl.GetAuthUser(r)
	// Return an error if the user is not authenticated
	if currentUser == nil {
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Parse the form data
	err := r.ParseForm()
	// Return an error if there was an issue parsing the form data
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Get the post_id and comment_content from the form data
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("comment_content")

	// Convert the post_id to an integer
	postID, err := strconv.Atoi(postIDStr)
	// Return an error if the post_id is not a valid integer
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Get all the posts and find the corresponding post for the provided post_id
	posts, _ := ctl.Repo.PRepo.GetAllPosts()
	var post model.Post
	for _, pos := range posts {
		if pos.ID == postID {
			post = *pos
		}
	}

	// Set the post's username to the current user's username
	post.Username = currentUser.Username
	// Return an error if the comment_content is empty
	if content == "" {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Create a new comment with the current user's information, post_id, and comment_content
	comment := model.Comment{
		UserID:    currentUser.ID,
		PostID:    postID,
		Content:   content,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now(),
	}
	// Create the comment in the repository
	ctl.Repo.ComRepo.CreateComment(comment)

	// Redirect the user to the post's page
	http.Redirect(w, r, fmt.Sprintf("/post?post_id=%d", postID), http.StatusSeeOther)
}

// DeleteComment handles the HTTP request to delete a comment.
//
// It expects a POST request, and checks if the current user is authenticated.
// It then parses the request form, retrieves the comment ID, and checks if the comment exists.
// If the current user is the owner of the comment, it deletes the comment and redirects the user to the post page.
func (ctl *BaseController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		// Return a 405 Method Not Allowed error if the method is not POST
		ctl.ErrorController(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Get the authenticated user
	currentUser := ctl.GetAuthUser(r)
	// Check if the user is authenticated
	if currentUser == nil {
		// Return a 401 Unauthorized error if the user is not authenticated
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	// Parse the request form
	err := r.ParseForm()
	if err != nil {
		// Return a 400 Bad Request error if the request form cannot be parsed
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Get the comment ID from the request form
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		// Return a 400 Bad Request error if the comment ID cannot be parsed
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Get the comment by ID
	comment, _ := ctl.Repo.ComRepo.GetCommentByID(commentID)
	// Check if the comment exists
	if comment == nil {
		// Return a 404 Not Found error if the comment does not exist
		ctl.ErrorController(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	// Check if the current user is the owner of the comment
	if !ctl.IsCommentOwner(currentUser, comment) {
		// Return a 403 Forbidden error if the user is not the owner of the comment
		ctl.ErrorController(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	// Get the post ID from the comment
	var postID int
	postID = comment.PostID

	// Delete the comment
	ctl.Repo.ComRepo.DeleteComment(commentID)
	// Redirect the user to the post page
	http.Redirect(w, r, fmt.Sprintf("/post?post_id=%d", postID), http.StatusSeeOther)
}

func (ctl *BaseController) IsCommentOwner(user *model.User, comment *model.Comment) bool {
	return user.ID == comment.UserID
}
