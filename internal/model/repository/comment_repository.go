package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"forum/internal/model"
	"forum/pkg/sqlite"
)

// CommentRepository
type CommentRepository struct {
	DB *sqlite.Database
}

// NewCommentRepository
func NewCommentRepository(db *sqlite.Database) *CommentRepository {
	return &CommentRepository{
		DB: db,
	}
}

// CreateComment
func (comr *CommentRepository) CreateComment(comment model.Comment) error {
	_, err := comr.DB.SQLite.Exec("INSERT INTO comments (post_id, user_id, content, likes, dislikes, created_at) VALUES(?, ?, ?, ?, ?, ?);",
		comment.PostID, comment.UserID, comment.Content, comment.Likes, comment.Dislikes, comment.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// GetCommentsByPostID
func (comr *CommentRepository) GetCommentsByPostID(postID int) ([]*model.Comment, error) {
	comRows, err := comr.DB.SQLite.Query("SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.likes, comments.dislikes, comments.created_at, users.username FROM comments LEFT JOIN users ON comments.user_id = users.id WHERE post_id = ?", postID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer comRows.Close()

	coms := []*model.Comment{}
	defer comRows.Close()
	for comRows.Next() {
		com := model.Comment{}
		if err := comRows.Scan(&com.ID, &com.PostID, &com.UserID, &com.Content, &com.Likes, &com.Dislikes, &com.CreatedAt, &com.Username); err != nil {
			slog.Error(err.Error())
			continue
		}
		coms = append(coms, &com)
	}

	return coms, nil
}

func (comr *CommentRepository) GetAllCommentData(commentId int) *model.Comment {
	com := model.Comment{}
	err := comr.DB.SQLite.QueryRow("SELECT id, post_id, user_id, content, likes, dislikes, created_at FROM comments WHERE id = ?", commentId).Scan(&com.ID, &com.PostID, &com.UserID, &com.Content, &com.Likes, &com.Dislikes, &com.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error(err.Error())
		}
		return nil
	}
	return &com
}

// GetCommentsByUserID
func (comr *CommentRepository) GetCommentsByUserID(userId int) ([]*model.Comment, error) {
	comRows, err := comr.DB.SQLite.Query("SELECT id, post_id, user_id, content, likes, dislikes, created_at FROM comments WHERE user_id = ?", userId)
	if err != nil {
		fmt.Println("Getting commentsByUserId failed: " + err.Error())
		slog.Error(err.Error())
		return nil, err
	}
	comments := make([]*model.Comment, 0)
	defer comRows.Close()

	for comRows.Next() {
		comment := model.Comment{}
		if err := comRows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.CreatedAt); err != nil {
			fmt.Println(err)
			slog.Error(err.Error())
			continue
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (comr *CommentRepository) GetCommentLikesByUserID(userID, commentID int) model.TotalLikesComment {
	totlike := model.TotalLikesComment{}
	err := comr.DB.SQLite.QueryRow("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&totlike.CommentID, &totlike.UserID)
	if err != nil {
		return model.TotalLikesComment{}
	}
	return totlike
}

func (comr *CommentRepository) GetCommentDislikesByUserID(userID, commentID int) model.TotalDislikesComment {
	var totdis model.TotalDislikesComment
	err := comr.DB.SQLite.QueryRow("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&totdis.CommentID, &totdis.UserID)
	if err != nil {
		return model.TotalDislikesComment{}
	}
	return totdis
}

func (comr *CommentRepository) GetAllCommentLikesByUserID(userID int) []*model.TotalLikesComment {
	totlike := []*model.TotalLikesComment{}
	tabl, err := comr.DB.SQLite.Query("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllCommentLikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		row := model.TotalLikesComment{}
		if err := tabl.Scan(&row.CommentID, &row.UserID); err != nil {
			return nil
		}
		totlike = append(totlike, &row)
	}
	return totlike
}

func (comr *CommentRepository) GetAllCommentDislikesByUserID(userID int) []*model.TotalDislikesComment {
	dis := []*model.TotalDislikesComment{}
	tabl, err := comr.DB.SQLite.Query("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("GetAllCommentDislikesByUserID failed: " + err.Error())
	}
	defer tabl.Close()
	for tabl.Next() {
		row := model.TotalDislikesComment{}
		if err := tabl.Scan(&row.CommentID, &row.UserID); err != nil {
			return nil
		}
		dis = append(dis, &row)
	}
	return dis
}

// GetCommentByID
func (comr *CommentRepository) GetCommentByID(commentId int) (*model.Comment, error) {
	com := model.Comment{}
	err := comr.DB.SQLite.QueryRow("SELECT id, post_id, user_id, content, likes, dislikes, created_at FROM comments WHERE id = ?", commentId).Scan(&com.ID, &com.PostID, &com.UserID, &com.Content, &com.Likes, &com.Dislikes, &com.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			slog.Error(err.Error())
		}
		return nil, err
	}
	return &com, nil
}

// UpdateComment
func (comr *CommentRepository) UpdateCommentData(commentID, newLikes, newDislikes int) error {
	_, err := comr.DB.SQLite.Exec("UPDATE comments SET likes = ?, dislikes = ? WHERE id = ?", newLikes, newDislikes, commentID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// DeleteComment
func (comr *CommentRepository) DeleteComment(commentID int) error {
	_, err := comr.DB.SQLite.Exec("DELETE FROM comments WHERE id = ?", commentID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
