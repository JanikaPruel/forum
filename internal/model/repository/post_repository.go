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

// Posts
func (pr *PostRepository) CreatePost(post model.Post) (categoryID int, err error) {
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
func (pr *PostRepository) GetAllPost() (posts []*model.Post, err error) {

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

// GetPostByID
func (pr *PostRepository) GetPostBy(postId int) (*model.Post, error) {
	post := model.Post{}

	err := pr.DB.SQLite.QueryRow("SELECT p.id, p.user_id, p.title, p.content, p.likes, p.dislikes, p.created_at, u.username FROM posts AS p LEFT JOIN users AS u ON p.user_id = u.id WHERE p.id = ?",
		postId).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt, &post.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error(err.Error())
		}
		return nil, err
	}
	return &post, nil
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
