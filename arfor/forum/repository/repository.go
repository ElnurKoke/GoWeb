package repository

import (
	"database/sql"
)

type Repository struct {
	Authorization
	Post
	Comment
	Like
	Dislike
}

type LikeDislike interface {
	// Post - Like//
	CheckPostLike(int, int) error
	DeletePostLike(int, int) error
	// Post - Dislike
	CheckPostDislike(int, int) error
	DeletePostDislike(int, int) error

	// Comment - Like

	CheckCommentLike(int, int) error
	DeleteCommentLike(int, int) error

	CheckCommentDislike(int, int) error
	DeleteCommentDislike(int, int) error
	// UpdateCommentVote(int) error
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
		Post:          NewPostRepo(db),
		Comment:       NewCommentRepo(db),
		Like:          NewLikeRepo(db),
		Dislike:       NewDislikeRepo(db),
	}
}
