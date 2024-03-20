package controller

import (
	"net/http"
	"sort"
	"strconv"

	"forum/internal/model"
)

func (ctl *BaseController) FilterByLikes(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Likes > posts[j].Likes
	})
}

func (ctl *BaseController) FilterByDislikes(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Dislikes > posts[j].Dislikes
	})
}

// AddLikePost handles the request to like a post.
// It takes in a http.ResponseWriter and a http.Request as parameters.
func (ctl *BaseController) AddLikeInPost(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("post_id")
	// Convert the post ID to an integer
	postID, _ := strconv.Atoi(postIDStr)
	// Get the authenticated user
	authUser := ctl.GetAuthUser(r)
	// Check if the user is authenticated
	if authUser == nil {
		// Return an unauthorized error if the user is not authenticated
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	// Add a reaction (like) to the post for the authenticated user
	ctl.AddReactionInPost(Reaction{ID: postID, UserID: authUser.ID, MarkID: "posts", Mark: "like"})
	// Set the content type of the response to "text/javascript"
	w.Header().Set("Content-Type", "text/javascript")
	// Write "window.history.back();" to the response body
	w.Write([]byte("window.history.back();"))
}

// HandlerDislikePost handles the request to dislike a post.
// It takes in a http.ResponseWriter and a http.Request as parameters.
// The function writes "window.history.back();" to the response body and sets the content type to "text/javascript".
func (ctl *BaseController) AddDislikeInPost(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("post_id")
	// Convert the post ID to an integer
	postID, _ := strconv.Atoi(postIDStr)
	// Get the authenticated user
	authUser := ctl.GetAuthUser(r)
	// Check if the user is authenticated
	if authUser == nil {
		// Return an unauthorized error if the user is not authenticated
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	// Add a reaction (dislike) to the post for the authenticated user
	ctl.AddReactionInPost(Reaction{ID: postID, UserID: authUser.ID, MarkID: "posts", Mark: "dislike"})
	// Set the content type of the response to "text/javascript"
	w.Header().Set("Content-Type", "text/javascript")
	// Write "window.history.back();" to the response body
	w.Write([]byte("window.history.back();"))
}

// AddReactionInPost adds a reaction to a post.
func (ctl *BaseController) AddReactionInPost(r Reaction) {
	// Retrieve the likes and dislikes of the user for the post
	likes := ctl.Repo.RRepo.GetPostLikesByUserID(r.UserID, r.ID)
	dislikes := ctl.Repo.RRepo.GetPostDislikesByUserID(r.UserID, r.ID)

	// Retrieve the post and user information
	post := &model.Post{}
	user, _ := ctl.Repo.URepo.GetUserByID(r.UserID)
	posts, _ := ctl.Repo.PRepo.GetAllPosts()

	// Find the post with the given ID
	for _, pos := range posts {
		if pos.ID == r.ID {
			post = pos
			break
		}
	}
	post.Username = user.Username

	// Add or remove reaction based on the type of reaction
	if r.Mark == "like" {
		// Remove the like if it already exists
		if likes != (model.TotalLikesPost{}) {
			ctl.Repo.PRepo.UpdatePostByID(r.ID, post.Likes-1, post.Dislikes, post.Comment)
			ctl.Repo.RRepo.RemovePostLikesByUserID(r.UserID, r.ID)
			return
		}
		// Decrease dislikes if it exists
		if dislikes != (model.TotalDislikesPost{}) {
			post.Dislikes--
		}
		ctl.Repo.PRepo.UpdatePostByID(r.ID, post.Likes+1, post.Dislikes, post.Comment)
		ctl.Repo.RRepo.RemovePostDislikesByUserID(r.UserID, r.ID)
		ctl.Repo.RRepo.InsertPostLike(model.TotalLikesPost{UserID: r.UserID, PostID: r.ID})
	} else if r.Mark == "dislike" {
		// Remove the dislike if it already exists
		if dislikes != (model.TotalDislikesPost{}) {
			ctl.Repo.PRepo.UpdatePostByID(r.ID, post.Likes, post.Dislikes-1, post.Comment)
			ctl.Repo.RRepo.RemovePostDislikesByUserID(r.UserID, r.ID)
			return
		}
		// Decrease likes if it exists
		if likes != (model.TotalLikesPost{}) {
			post.Likes--
		}
		ctl.Repo.PRepo.UpdatePostByID(r.ID, post.Likes, post.Dislikes+1, post.Comment)
		ctl.Repo.RRepo.RemovePostLikesByUserID(r.UserID, r.ID)
		ctl.Repo.RRepo.InsertPostDislike(model.TotalLikesPost{UserID: r.UserID, PostID: r.ID})
	}
}

func (ctl *BaseController) AddLikeInComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Query().Get("comment_id")
	commentID, _ := strconv.Atoi(commentIDStr)
	authUser := ctl.GetAuthUser(r)
	if authUser == nil {
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	ctl.AddReactionInComment(Reaction{ID: commentID, UserID: authUser.ID, MarkID: "comments", Mark: "like"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

func (ctl *BaseController) AddDislikeInComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.URL.Query().Get("comment_id")
	commentID, _ := strconv.Atoi(commentIDStr)
	authUser := ctl.GetAuthUser(r)
	if authUser == nil {
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	ctl.AddReactionInComment(Reaction{ID: commentID, UserID: authUser.ID, MarkID: "comments", Mark: "dislike"})
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte("window.history.back();"))
}

// AddReactionInComment updates comment data based on the reaction received.
// It handles both likes and dislikes reactions.
//
// Parameters:
// - r: The reaction data containing the user ID, comment ID, and reaction type (like or dislike).
func (ctl *BaseController) AddReactionInComment(r Reaction) {
	// Get the likes and dislikes for the user and comment
	// from the comment repository
	likes := ctl.Repo.ComRepo.GetCommentLikesByUserID(r.UserID, r.ID)
	dislikes := ctl.Repo.ComRepo.GetCommentDislikesByUserID(r.UserID, r.ID)

	// Get all the comment data from the comment repository
	comments := *ctl.Repo.ComRepo.GetAllCommentData(r.ID)

	// Handle likes reaction
	if r.Mark == "like" {
		// If the user already liked the comment,
		// remove the like and return
		if likes != (model.TotalLikesComment{}) {
			ctl.Repo.ComRepo.UpdateCommentData(r.ID, comments.Likes-1, comments.Dislikes)
			ctl.Repo.RRepo.RemoveCommentLikesByUserID(r.UserID, r.ID)
			return
		}
		// If the user already disliked the comment,
		// subtract one dislike and add one like
		if dislikes != (model.TotalDislikesComment{}) {
			comments.Dislikes--
		}
		ctl.Repo.ComRepo.UpdateCommentData(r.ID, comments.Likes+1, comments.Dislikes)
		ctl.Repo.RRepo.RemoveCommentDislikesByUserID(r.UserID, r.ID)
		ctl.Repo.RRepo.InsertCommentLike(model.TotalLikesComment{UserID: r.UserID, CommentID: r.ID})

	} else if r.Mark == "dislike" {
		// If the user already disliked the comment,
		// remove the dislike and return
		if dislikes != (model.TotalDislikesComment{}) {
			ctl.Repo.ComRepo.UpdateCommentData(r.ID, comments.Likes, comments.Dislikes-1)
			ctl.Repo.RRepo.RemoveCommentDislikesByUserID(r.UserID, r.ID)
			return
		}
		// If the user already liked the comment,
		// subtract one like and add one dislike
		if likes != (model.TotalLikesComment{}) {
			comments.Likes--
		}
		ctl.Repo.ComRepo.UpdateCommentData(r.ID, comments.Likes, comments.Dislikes+1)
		ctl.Repo.RRepo.RemoveCommentLikesByUserID(r.UserID, r.ID)
		ctl.Repo.RRepo.InsertCommentDislike(model.TotalDislikesComment{UserID: r.UserID, CommentID: r.ID})
	}
}
