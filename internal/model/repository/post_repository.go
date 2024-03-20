package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"forum/internal/model"
	"forum/pkg/sqlite"
)

// PostRepository
type PostRepository struct {
	DB *sqlite.Database
}

// NewPostRepository
func NewPostRepository(db *sqlite.Database) *PostRepository {
	return &PostRepository{
		DB: db,
	}
}

// CreatePost
func (pr *PostRepository) CreatePost(post model.Post) (ID int, err error) {
	res, err := pr.DB.SQLite.Exec("INSERT INTO posts (user_id, title, content, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		post.UserID, post.Title, post.Content, post.Likes, post.Dislikes, post.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	for _, categoryID := range post.Category {
		_, err = pr.DB.SQLite.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", id, categoryID)
		if err != nil {
			slog.Error(err.Error())
			return 0, err
		}
	}
	return int(id), nil
}

// GetAllPosts
func (pr *PostRepository) GetAllPosts() (posts []*model.Post, err error) {

	postsRows, err := pr.DB.SQLite.Query("SELECT id, user_id, title, content, likes, dislikes, created_at FROM posts")
	if err != nil {
		slog.Error(err.Error())

	}
	defer postsRows.Close()

	for postsRows.Next() {
		post := model.Post{}

		if err := postsRows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}
		categoryRows, err := pr.DB.SQLite.Query("SELECT category_id FROM post_categories WHERE post_id = ?", post.ID)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		defer categoryRows.Close()

		categories := []int{}
		for categoryRows.Next() {
			var categoryID int

			if err := categoryRows.Scan(&categoryID); err != nil {
				fmt.Println("Error scanning category ID: " + err.Error())
				continue
			}
			categories = append(categories, categoryID)
		}
		post.Category = categories

		posts = append(posts, &post)
	}

	return posts, nil
}

// GetPostByUserID
func (pr *PostRepository) GetPostsByUserID(userID int) (posts []*model.Post, err error) {
	postsRows, err := pr.DB.SQLite.Query("SELECT id, user_id, title, content, likes, dislikes, created_at FROM posts WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer postsRows.Close()

	for postsRows.Next() {
		post := model.Post{}

		if err := postsRows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (pr *PostRepository) GetAllPostLikesByUserID(userID int) []*model.TotalLikesPost {
	like := []*model.TotalLikesPost{}
	tabl, err := pr.DB.SQLite.Query("SELECT post_id, user_id FROM total_likes_post WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		row := model.TotalLikesPost{}
		if err := tabl.Scan(&row.PostID, &row.UserID); err != nil {
			return nil
		}
		like = append(like, &row)
	}
	return like
}

func (pr *PostRepository) GetAllPostDislikesByUserID(userID, postID int) model.TotalDislikesPost {
	dis := model.TotalDislikesPost{}
	err := pr.DB.SQLite.QueryRow("SELECT post_id, user_id FROM total_dislikes_post WHERE user_id = ? AND post_id = ?", userID, postID).
		Scan(&dis.PostID, &dis.UserID)
	if err != nil {
		return model.TotalDislikesPost{}
	}
	return dis
}

func (pr *PostRepository) GetAllPostDislikesByUserIDs(userID int) []*model.TotalDislikesPost {
	dis := []*model.TotalDislikesPost{}
	tabl, err := pr.DB.SQLite.Query("SELECT post_id, user_id FROM total_dislikes_post WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		row := model.TotalDislikesPost{}
		if err := tabl.Scan(&row.PostID, &row.UserID); err != nil {
			return nil
		}
		dis = append(dis, &row)
	}
	return dis
}

// GetPostByID
func (pr *PostRepository) GetPostByID(postID int) *model.Post {
	row, err := pr.DB.SQLite.Query("SELECT * FROM posts WHERE id=?", postID)
	// row, err := pr.DB.SQLite.Query("SELECT p.id, ///p.user_id, p.title, p.content, p.likes, p.dislikes, p.created_at, u.username posts AS p LEFT JOIN users AS u ON p.user_id = u.id WHERE p.id = ?", postID)
	// err := pr.DB.SQLite.QueryRow("SELECT id, user_id, title, content, likes, dislikes FROM posts WHERE id=?", postID).
	// Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes)
	defer row.Close()
	if err == nil {
		for row.Next() {
			p := model.Post{}
			err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.Username)
			if err != nil {
				fmt.Println("ERROR DB, err:", err.Error())
				continue
			}
			if p.ID == postID {
				return &p
			}
		}
	}
	return nil
}

func (pr *PostRepository) GetAllPostData(postId int) *model.Post {
	var p model.Post
	err := pr.DB.SQLite.QueryRow("SELECT p.id, p.user_id, p.title, p.content, p.likes, p.dislikes, p.created_at, u.username FROM posts AS p LEFT JOIN users AS u ON p.user_id = u.id WHERE p.id = ?", postId).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Likes, &p.Dislikes, &p.CreatedAt, &p.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get all post data failed: " + err.Error())
		}
		return nil
	}
	return &p
}

// UpdatePostByID
func (pr *PostRepository) UpdatePostByID(postId, likes, dislikes, comments int) error {
	_, err := pr.DB.SQLite.Exec("UPDATE posts SET likes = ?, dislikes = ? WHERE id = ?;", likes, dislikes, postId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// DeletePostByID
func (pr *PostRepository) DeletePost(postID int) error {
	_, err := pr.DB.SQLite.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		fmt.Println("Deleting post failed: " + err.Error())
	}

	return nil
}
