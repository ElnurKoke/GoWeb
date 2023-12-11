package repository

import (
	"database/sql"
	"fmt"
	"log"

	"forum/models"
)

type LikeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) *LikeRepo {
	return &LikeRepo{
		db: db,
	}
}

func (r *LikeRepo) SetPostLike(like models.Like) error {
	query := `INSERT INTO like(user_id, post_id, active) VALUES($1, $2, $3)`
	_, err := r.db.Exec(query, like.UserID, like.PostID, 1)
	if err != nil {
		return fmt.Errorf(path+"set post like: %w", err)
	}
	return nil
}

func (r *LikeRepo) CheckPostLike(userID, postID int) error {
	query := `SELECT id FROM like WHERE user_id = $1 AND post_id = $2 AND active = 1`
	row := r.db.QueryRow(query, userID, postID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (r *LikeRepo) CheckPostDislike(userID, postID int) error {
	query := `SELECT id FROM dislike WHERE user_id = $1 AND post_id = $2 AND active = 1`
	row := r.db.QueryRow(query, userID, postID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (r *LikeRepo) DeletePostLike(userID, postID int) error {
	query := `DELETE FROM like WHERE user_id = $1 AND post_id = $2`
	_, err := r.db.Exec(query, userID, postID)
	if err != nil {
		return fmt.Errorf(path+"delete post like: %w", err)
	}
	return nil
}

func (r *LikeRepo) DeletePostDislike(userID, postID int) error {
	query := `DELETE FROM dislike WHERE user_id = $1 AND post_id = $2`
	_, err := r.db.Exec(query, userID, postID)
	if err != nil {
		return fmt.Errorf(path+"delete post dislike: %w", err)
	}
	return nil
}

func (r *LikeRepo) UpdatePostVote(postID int) error {
	query := `SELECT COUNT(post_id) FROM like WHERE post_id = $1 AND active = $2`
	row := r.db.QueryRow(query, postID, 1)
	var likesCount int
	if err := row.Scan(&likesCount); err != nil {
		return fmt.Errorf(path+"update post like: scan like: %w", err)
	}
	query2 := `SELECT COUNT(post_id) FROM dislike WHERE post_id = $1 AND active = $2`
	row = r.db.QueryRow(query2, postID, 1)
	var dislikesCount int
	if err := row.Scan(&dislikesCount); err != nil {
		return fmt.Errorf(path+"update post like: scan dislike: %w", err)
	}
	query3 := `UPDATE post SET like = $1, dislike = $2 WHERE id = $3`
	_, err := r.db.Exec(query3, likesCount, dislikesCount, postID)
	if err != nil {
		return fmt.Errorf(path+"update post like: exec: %w", err)
	}
	return nil
}

//------------------Commnet---------------------------//

func (r *LikeRepo) SetCommentLike(like models.Like) error {
	query := `INSERT INTO like(user_id, comment_id, active) VALUES($1, $2, $3)`
	_, err := r.db.Exec(query, like.UserID, like.CommentID, 1)
	if err != nil {
		log.Printf(path+"set comment like: %s", err)
		return fmt.Errorf(path+"set post like: %w", err)
	}
	return nil
}

func (r *LikeRepo) CheckCommentLike(userID, commentID int) error {
	query := `SELECT id FROM like WHERE user_id = $1 AND comment_id = $3 AND active = 1`
	row := r.db.QueryRow(query, userID, commentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		fmt.Printf("chech comment like: %s\n", err)
		return err
	}
	return nil
}

func (r *LikeRepo) CheckCommentDislike(userID, commentID int) error {
	query := `SELECT id FROM dislike WHERE user_id = $1 AND comment_id = $3 AND active = 1`
	row := r.db.QueryRow(query, userID, commentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		fmt.Printf("chech comment dislike: %s\n", err)
		return err
	}
	return nil
}

func (r *LikeRepo) DeleteCommentLike(userID, commentID int) error {
	query := `DELETE FROM like WHERE user_id = $1 AND comment_id = $2`
	_, err := r.db.Exec(query, userID, commentID)
	if err != nil {
		log.Printf(path+"delete comment like: %s", err)
		return fmt.Errorf(path+"delete post like: %w", err)
	}
	return nil
}

func (r *LikeRepo) DeleteCommentDislike(userID, commentID int) error {
	query := `DELETE FROM dislike WHERE user_id = $1 AND comment_id = $2`
	_, err := r.db.Exec(query, userID, commentID)
	if err != nil {
		log.Printf(path+"delete comment dislike: %s", err)
		return fmt.Errorf(path+"delete post dislike: %w", err)
	}
	return nil
}

func (r *LikeRepo) UpdateCommentVote(commentID int) (int, error) {
	query := `SELECT COUNT(id) FROM like WHERE comment_id = $1 AND active = $2`
	row := r.db.QueryRow(query, commentID, 1)
	var likesCount int
	if err := row.Scan(&likesCount); err != nil {
		fmt.Printf("scan comment like count: %s\n", err)
		return 0, fmt.Errorf(path+"update post like: scan like: %w", err)
	}
	query2 := `SELECT COUNT(id) FROM dislike WHERE comment_id = $1 AND active = $2`
	row = r.db.QueryRow(query2, commentID, 1)
	var dislikesCount int
	if err := row.Scan(&dislikesCount); err != nil {
		fmt.Printf("scan comment dislike count: %s\n", err)
		return 0, fmt.Errorf(path+"update post like: scan dislike: %w", err)
	}
	query3 := `UPDATE comment SET like = $1, dislike = $2 WHERE id = $3`
	_, err := r.db.Exec(query3, likesCount, dislikesCount, commentID)
	if err != nil {
		fmt.Printf("update comment vote: %s\n", err)
		return 0, fmt.Errorf(path+"update post like: exec: %w", err)
	}
	var postID int
	query4 := `SELECT post_id FROM comment WHERE id = $1`
	row = r.db.QueryRow(query4, commentID)
	if err := row.Scan(&postID); err != nil {
		return 0, fmt.Errorf(path+"update comment vote: select post_id: %w", err)
	}
	return postID, nil
}
