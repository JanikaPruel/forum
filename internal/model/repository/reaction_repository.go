package repository

import (
	"log/slog"

	"forum/internal/model"
	"forum/pkg/sqlite"
)

// ReactionRepository (likes, and dislikes)
type ReactionRepository struct {
	DB *sqlite.Database
}

// NewReactionRepository
func NewReactionRepository(db *sqlite.Database) *ReactionRepository {
	return &ReactionRepository{
		DB: db,
	}
}

// InsertCommentLike
func (re *ReactionRepository) InsertCommentLike(commentLike model.TotalLikesComment) error {
	_, err := re.DB.SQLite.Exec("INSERT OR IGNORE INTO total_likes_comment (comment_id, user_id) VALUES (?, ?)", commentLike.CommentID, commentLike.UserID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

// InsertCommentDislike
func (re *ReactionRepository) InsertCommentDislike(commentDislike model.TotalDislikesComment) error {
	_, err := re.DB.SQLite.Exec("INSERT OR IGNORE INTO total_dislikes_comment (comment_id, user_id) VALUES (?, ?)", commentDislike.CommentID, commentDislike.UserID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// GetCommentLikesByUserID
func (re *ReactionRepository) GetCommentLikesByUserID(userID, commentID int) (totlike model.TotalLikesComment, err error) {
	err = re.DB.SQLite.QueryRow("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&totlike.CommentID, &totlike.UserID)
	if err != nil {
		slog.Error(err.Error())
		return totlike, err
	}
	return totlike, nil
}

// GetCommentDislikesByUserID
func (re *ReactionRepository) GetCommentDislikesByUserID(userID, commentID int) (totdis model.TotalLikesComment, err error) {
	err = re.DB.SQLite.QueryRow("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID).
		Scan(&totdis.CommentID, &totdis.UserID)
	if err != nil {
		slog.Error(err.Error())
		return totdis, err
	}
	return totdis, nil
}

// GetAllCommentLikesByUserID
func (re *ReactionRepository) GetAllCommentLikesByUserID(userID int) ([]*model.TotalLikesComment, error) {
	totlike := []*model.TotalLikesComment{}
	rows, err := re.DB.SQLite.Query("SELECT comment_id, user_id FROM total_likes_comment WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		like := model.TotalLikesComment{}
		if err := rows.Scan(&like.CommentID, &like.UserID); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		totlike = append(totlike, &like)
	}
	return totlike, nil
}

// GetAllCommentDislikesByUserID
func (re *ReactionRepository) GetAllCommentDislikesByUserID(userID int) ([]*model.TotalDislikesComment, error) {
	totdis := []*model.TotalDislikesComment{}
	rows, err := re.DB.SQLite.Query("SELECT comment_id, user_id FROM total_dislikes_comment WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dis := model.TotalDislikesComment{}
		if err := rows.Scan(&dis.CommentID, &dis.UserID); err != nil {
			return nil, err
		}
		totdis = append(totdis, &dis)
	}
	return totdis, nil
}

// GetLikedCommentsByUserID
func (re *ReactionRepository) GetLikedCommentsByUserID(userID int) ([]*model.Comment, error) {
	rows, err := re.DB.SQLite.Query("SELECT c.id, c.user_id, c.post_id, c.content, c.likes, c.dislikes FROM comments c JOIN total_likes_comment lc ON c.id = lc.comment_id WHERE lc.user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		com := model.Comment{}
		err = rows.Scan(&com.ID, &com.UserID, &com.PostID, &com.Content, &com.Likes, &com.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &com)
	}
	return comments, nil
}

// RemoveCommentLikesByUserID
func (re *ReactionRepository) RemoveCommentLikesByUserID(userID, commentID int) error {
	_, err := re.DB.SQLite.Exec("DELETE FROM total_likes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		return err
	}
	return nil
}

// RemoveCommentDislikesByUserID
func (re *ReactionRepository) RemoveCommentDislikesByUserID(userID, commentID int) error {
	_, err := re.DB.SQLite.Exec("DELETE FROM total_dislikes_comment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		return err
	}
	return nil
}
