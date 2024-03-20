package controller

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"forum/internal/model"
)

type PostData struct {
	AuntUser        *model.User
	Categories      []*model.Category
	Post            *model.Post
	Comments        *model.Comment
	PostLikes       *model.TotalLikesPost
	PostDislikes    *model.TotalDislikesPost
	CommentLikes    *model.TotalLikesComment
	CommentDislikes *model.TotalDislikesComment
}

// CreataPost | POST /posts
func (ctl *BaseController) CreatePost(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(GetWD() + "/web/templates/create_post.html"))

	user := ctl.GetAuthUser(r)
	logged := false
	if user != nil {
		logged = true
	}

	// categories
	categories, err := ctl.Repo.CRepo.GetAllCategories()
	if err != nil {
		slog.Error(err.Error())
	}

	err = r.ParseForm()
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	title := r.FormValue("post-title")
	content := r.FormValue("post-content")
	categoryIDsStr := r.Form["category-create-post[]"]
	if len(categoryIDsStr) == 0 {
		data := struct {
			Err        int
			Message    string
			Categories []*model.Category
			UserID     *model.User
			IsLoggedIn bool
		}{
			Err:        1,
			Message:    "Choose at least one category",
			Categories: categories,
			UserID:     user,
			IsLoggedIn: logged,
		}
		tmpl.Execute(w, data)
		return
	}

	categoryIDs := []int{}
	for _, categoryIDStr := range categoryIDsStr {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			data := struct {
				Err        int
				Message    string
				Categories []*model.Category
				UserID     *model.User
				IsLoggedIn bool
			}{
				Err:        1,
				Message:    "Choose at least one category",
				Categories: categories,
				UserID:     user,
				IsLoggedIn: logged,
			}
			tmpl.Execute(w, data)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	currentUser := ctl.GetAuthUser(r)

	post := model.Post{
		UserID:    currentUser.ID,
		Title:     title,
		Content:   content,
		Category:  categoryIDs,
		Likes:     0,
		Dislikes:  0,
		CreatedAt: time.Now(),
	}
	if title == "" || content == "" {
		err := 4
		mes := "Title and content must be not empty, enter your data, and try again!"
		data := struct {
			Err     int
			Message string
		}{
			Err:     err,
			Message: mes,
		}
		tmpl.Execute(w, data)
		return
	}
	_, err = ctl.Repo.PRepo.CreatePost(post)
	if err != nil {
		slog.Error(err.Error())
		data := struct {
			Err        int
			Message    string
			Categories []*model.Category
			UserID     *model.User
			IsLoggedIn bool
		}{
			Err:        1,
			Message:    "Invalid data, try again!",
			Categories: categories,
			UserID:     user,
			IsLoggedIn: logged,
		}
		tmpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, "GET /", http.StatusSeeOther)
}

func (ctl *BaseController) SearchPosts(posts []*model.Post, searchQuery string, categoryFilters []int) []*model.Post {
	var filteredPosts []*model.Post

	for _, post := range posts {
		if len(categoryFilters) > 0 && !contains(post.Category, categoryFilters[0]) {
			continue
		}
		if strings.Contains(strings.ToLower(post.Title), strings.ToLower(searchQuery)) ||
			strings.Contains(strings.ToLower(post.Content), strings.ToLower(searchQuery)) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

func (ctl *BaseController) GeneratePreviewContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	return content[:maxChars] + "..."
}

func contains(slice []int, item int) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func (ctl *BaseController) SortPostsByDate(posts []*model.Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
}

// ViewPostByID
func (ctl *BaseController) ViewPostByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("VIEW POST")
	if !strings.HasPrefix(r.URL.Path, "/post") {
		ctl.ErrorController(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		fmt.Println("THIS HERE VIEW")
		return
	}

	if r.Method != "GET" {
		ctl.ErrorController(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	postID := r.URL.Query().Get("post_id")
	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		ctl.ErrorController(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// post := ctl.Repo.PRepo.GetPostByID(postIdInt)
	post := &model.Post{}
	posts, _ := ctl.Repo.PRepo.GetAllPosts()
	for _, pos := range posts {
		if postIdInt == pos.ID {
			post = pos
		}
	}

	use, _ := ctl.Repo.URepo.GetUserByID(post.UserID)
	post.Username = use.Username

	postCategories, _ := ctl.Repo.CRepo.GetCategoriesByPostID(post.ID)

	postData := PostData{}
	postData.Categories = postCategories

	authUser := ctl.GetAuthUser(r)
	if authUser == nil {
		authUser = &model.User{}
	}
	postData.AuntUser = authUser

	postData.Post = post
	postLikes := ctl.Repo.PRepo.GetAllPostLikesByUserID(postData.AuntUser.ID)
	if len(postLikes) > 0 {
		postData.PostLikes = postLikes[0]
	}

	postDislikes := ctl.Repo.PRepo.GetAllPostDislikesByUserID(postData.AuntUser.ID, postData.Post.ID)

	postDislikesList := []*model.TotalDislikesPost{&postDislikes}
	if len(postDislikesList) > 0 && postDislikesList[0] != nil {
		postData.PostDislikes = postDislikesList[0]
	}

	comments, _ := ctl.Repo.ComRepo.GetCommentsByPostID(postData.Post.ID)
	ctl.SortCommentsByDate(comments)

	for _, c := range comments {
		commentLike := ctl.Repo.ComRepo.GetCommentLikesByUserID(postData.AuntUser.ID, c.ID)
		commentDislike := ctl.Repo.ComRepo.GetCommentDislikesByUserID(postData.AuntUser.ID, c.ID)
		if commentLike.CommentID != 0 {
			c.IsLikedByAuthUser = true
		}

		if commentDislike.CommentID != 0 {
			c.IsDislikedByAuthUser = true
		}
	}

	user := ctl.GetAuthUser(r)
	logged := false

	if user != nil {
		logged = true
	}

	likedPosts := []*model.TotalLikesPost{}
	dislikedPosts := []*model.TotalDislikesPost{}
	if logged {
		likedPosts = ctl.Repo.PRepo.GetAllPostLikesByUserID(user.ID)
		dislikedPosts = ctl.Repo.PRepo.GetAllPostDislikesByUserIDs(user.ID)
	}

	isLiked := false
	for _, likePost := range likedPosts {
		if likePost.PostID == postData.Post.ID {
			isLiked = true
		}
	}
	postData.Post.IsLikedByAuthUser = isLiked

	isDisliked := false
	for _, dislikePost := range dislikedPosts {
		if dislikePost.PostID == postData.Post.ID {
			isDisliked = true
		}
	}
	postData.Post.IsDislikedByAuthUser = isDisliked

	type Posts_Comments struct {
		UserID     *model.User
		IsLoggedIn bool
		*model.Post
		Comments   []*model.Comment
		AuthUserID int
		Categories []*model.Category
	}

	p_c := &Posts_Comments{
		UserID:     user,
		IsLoggedIn: logged,
		Post:       postData.Post,
		Comments:   comments,
		AuthUserID: postData.AuntUser.ID,
		Categories: postCategories,
	}

	tmp := template.Must(template.ParseFiles(GetWD() + "/web/templates/post.html"))

	if err := tmp.Execute(w, p_c); err != nil {
		slog.Error(err.Error())
	}
}

func (ctl *BaseController) SortCommentsByDate(comments []*model.Comment) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.After(comments[j].CreatedAt)
	})
}
