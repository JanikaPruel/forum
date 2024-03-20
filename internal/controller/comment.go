package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"forum/internal/model"
)

// CreateComment
func (ctl *BaseController) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ctl.ErrorController(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	currentUser := ctl.GetAuthUser(r)
	if currentUser == nil {
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err := r.ParseForm()
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("comment_content")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	post := model.Post{}

	posts, _ := ctl.Repo.PRepo.GetAllPosts()
	for _, pos := range posts {
		if pos.ID == postID {
			post = *pos
		}
	}

	post.Username = currentUser.Username
	if content == "" {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	comment := model.Comment{
		UserID:    currentUser.ID,
		PostID:    postID,
		Content:   content,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now(),
	}
	ctl.Repo.ComRepo.CreateComment(comment)

	http.Redirect(w, r, fmt.Sprintf("/post?post_id=%d", postID), http.StatusSeeOther)
}

// DeleteComment
func (ctl *BaseController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ctl.ErrorController(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	currentUser := ctl.GetAuthUser(r)
	if currentUser == nil {
		ctl.ErrorController(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err := r.ParseForm()
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	comment, _ := ctl.Repo.ComRepo.GetCommentByID(commentID)
	if comment == nil {
		ctl.ErrorController(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if !ctl.IsCommentOwner(currentUser, comment) {
		ctl.ErrorController(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	var postID int
	postID = comment.PostID

	ctl.Repo.ComRepo.DeleteComment(commentID)
	http.Redirect(w, r, fmt.Sprintf("/post?post_id=%d", postID), http.StatusSeeOther)
}

func (ctl *BaseController) IsCommentOwner(user *model.User, comment *model.Comment) bool {
	return user.ID == comment.UserID
}
