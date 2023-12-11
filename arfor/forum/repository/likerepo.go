package repository

import (
	"database/sql"
	"fmt"
	"foruum/models"
)

type Like interface {
	SetPostLike(models.Like) error
	SetCommentLike(models.Like) error
	LikeDislike
}

type LikeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) *LikeRepo {
	return &LikeRepo{
		db: db,
	}
}

func (l *LikeRepo) SetPostLike(like models.Like) error {
	query := `INSERT INTO like(postId, userId, active) VALUES(?,?,?) `
	_, err := l.db.Exec(query, like.PostID, like.UserID, like.Active)
	if err != nil {
		fmt.Println("SETPOSTLIKE error" + err.Error())
	}
	query = `update post set like=(select count(*) from like where postId=?) where id=?`
	_, err = l.db.Exec(query, like.PostID, like.PostID)
	if err != nil {
		fmt.Println("SETPOSTLIKE error" + err.Error())
	}
	return nil
}

func (l *LikeRepo) CheckPostLike(PostID, UserID int) error {
	query := `select id from like where userId = ? and postId = ?`
	row := l.db.QueryRow(query, UserID, PostID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (l *LikeRepo) DeletePostLike(PostID, UserID int) error {
	query := `delete from like where userId = ? and postId = ?`
	_, err := l.db.Exec(query, UserID, PostID)
	if err != nil {
		return fmt.Errorf("delete post like: %w", err)
	}
	query = `update post set like=(select count(*) from like where postId=?) where id=?`
	_, err = l.db.Exec(query, PostID, PostID)
	if err != nil {
		fmt.Println("SETPOSTLIKE error" + err.Error())
	}
	return nil
}

func (d *LikeRepo) CheckPostDislike(PostID, UserID int) error {
	query := `select id from dislike where userId = ? and postId = ?`
	row := d.db.QueryRow(query, UserID, PostID)
	var dislikeID int
	if err := row.Scan(&dislikeID); err != nil {
		return err
	}
	return nil
}

func (d *LikeRepo) DeletePostDislike(PostID, UserID int) error {
	query := `delete from dislike where userId = ? and postId = ?`
	_, err := d.db.Exec(query, UserID, PostID)
	if err != nil {
		return fmt.Errorf("delete post dislike: %w", err)
	}
	query = `update post set dislike=(select count(*) from dislike where postId=?) where id=?`
	_, err = d.db.Exec(query, PostID, PostID)
	if err != nil {
		fmt.Println("SETPOSTDISLIKE error" + err.Error())
	}
	return nil
}

// comment
func (l *LikeRepo) SetCommentLike(like models.Like) error {
	query := `INSERT INTO like(commentId,userId,active) VALUES(?, ?, ?) `
	_, err := l.db.Exec(query, like.CommentId, like.UserID, like.Active)
	if err != nil {
		fmt.Println("SETCOMMENTLIKE error" + err.Error())
	}
	query = `update comment set like=(select count(*) from like where commentId=?) where id=?`
	_, err = l.db.Exec(query, like.CommentId, like.CommentId)
	if err != nil {
		fmt.Println("SETCOMMENTLIKE error" + err.Error())
	}
	return nil
}

func (l *LikeRepo) CheckCommentLike(CommentID, UserID int) error {

	query := `select id from like where userId = ? and commentId = ?`
	row := l.db.QueryRow(query, UserID, CommentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (l *LikeRepo) DeleteCommentLike(CommentID, UserID int) error {
	query := `delete from like where userId = ? and commentId = ?`
	_, err := l.db.Exec(query, UserID, CommentID)
	if err != nil {
		return fmt.Errorf("delete post like: %w", err)
	}
	query = `update comment set like=(select count(*) from like where commentId=?) where id=?`
	_, err = l.db.Exec(query, CommentID, CommentID)
	if err != nil {
		fmt.Println("SETCOMMENTLIKE error" + err.Error())
	}
	return nil
}

func (l *LikeRepo) CheckCommentDislike(CommentID, UserID int) error {
	query := `select id from dislike where userId = ? and commentId = ?`
	row := l.db.QueryRow(query, UserID, CommentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (d *LikeRepo) DeleteCommentDislike(CommentID, UserID int) error {
	query := `delete from dislike where userId = ? and commentId = ?`
	_, err := d.db.Exec(query, UserID, CommentID)
	if err != nil {
		return fmt.Errorf("delete post like: %w", err)
	}
	query = `update comment set dislike=(select count(*) from dislike where commentId=?) where id=?`
	_, err = d.db.Exec(query, CommentID, CommentID)
	if err != nil {
		fmt.Println("SET COMMENT DISLIKE error" + err.Error())
	}
	return nil
}
