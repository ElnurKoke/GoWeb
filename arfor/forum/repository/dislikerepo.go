package repository

import (
	"database/sql"
	"fmt"
	"foruum/models"
)

type Dislike interface {
	SetPostDislike(models.Dislike) error
	SetCommentDislike(models.Dislike) error
	LikeDislike
}

type DislikeRepo struct {
	db *sql.DB
}

func NewDislikeRepo(db *sql.DB) *DislikeRepo {
	return &DislikeRepo{
		db: db,
	}
}

func (d *DislikeRepo) SetPostDislike(dislike models.Dislike) error {
	query := `INSERT INTO dislike(postId, userId,active) VALUES(?,?,?) `
	_, err := d.db.Exec(query, dislike.PostID, dislike.UserID, dislike.Active)
	if err != nil {
		fmt.Printf("SetPostLike error:%v", err)
	}
	query = `update post set dislike=(select count(*) from dislike where postId=?) where id=?`
	_, err = d.db.Exec(query, dislike.PostID, dislike.PostID)
	if err != nil {
		fmt.Println("SETPOSTLIKE error" + err.Error())
	}
	return nil
}

func (d *DislikeRepo) CheckPostDislike(PostID, UserID int) error {
	query := `select id from dislike where userId = ? and postId = ?`
	row := d.db.QueryRow(query, UserID, PostID)
	var dislikeID int
	if err := row.Scan(&dislikeID); err != nil {
		return err
	}
	return nil
}

func (d *DislikeRepo) DeletePostDislike(PostID, UserID int) error {
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

func (l *DislikeRepo) CheckPostLike(PostID, UserID int) error {
	query := `select id from like where userId = ? and postId = ?`
	row := l.db.QueryRow(query, UserID, PostID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (l *DislikeRepo) DeletePostLike(PostID, UserID int) error {
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

// comment
func (d *DislikeRepo) SetCommentDislike(dislike models.Dislike) error {
	query := `INSERT INTO dislike(commentId,userId,active) VALUES(?, ?, ?) `
	_, err := d.db.Exec(query, dislike.CommentId, dislike.UserID, dislike.Active)
	if err != nil {
		fmt.Printf("SetCommentDislike error:%v", err)
	}
	query = `update comment set dislike=(select count(*) from dislike where commentId=?) where id=?`
	_, err = d.db.Exec(query, dislike.CommentId, dislike.CommentId)
	if err != nil {
		fmt.Println("SET COMMENT DISLIKE error" + err.Error())
	}
	return nil
}

func (l *DislikeRepo) CheckCommentDislike(CommentID, UserID int) error {
	query := `select id from dislike where userId = ? and commentId = ?`
	row := l.db.QueryRow(query, UserID, CommentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (d *DislikeRepo) DeleteCommentDislike(CommentID, UserID int) error {
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

func (l *DislikeRepo) CheckCommentLike(CommentID, UserID int) error {
	query := `select id from like where userId = ? and commentId = ?`
	row := l.db.QueryRow(query, UserID, CommentID)
	var likeID int
	if err := row.Scan(&likeID); err != nil {
		return err
	}
	return nil
}

func (l *DislikeRepo) DeleteCommentLike(CommentID, UserID int) error {
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
