package service

import (
	"foruum/models"
	"foruum/repository"
)

type LikeService struct {
	repo repository.Like
}

func NewLikeService(repo repository.Like) *LikeService {
	return &LikeService{
		repo: repo,
	}
}

func (l *LikeService) SetPostLike(like models.Like) error {
	if l.repo.CheckPostDislike(like.PostID, like.UserID) == nil {
		l.repo.DeletePostDislike(like.PostID, like.UserID)
	}
	if l.repo.CheckPostLike(like.PostID, like.UserID) != nil {
		return l.repo.SetPostLike(like)
	} else {
		return l.repo.DeletePostLike(like.PostID, like.UserID)
	}
}

func (l *LikeService) SetCommentLike(like models.Like) error {
	if l.repo.CheckCommentDislike(like.CommentId, like.UserID) == nil {
		l.repo.DeleteCommentDislike(like.CommentId, like.UserID)
	}
	if l.repo.CheckCommentLike(like.CommentId, like.UserID) != nil {
		return l.repo.SetCommentLike(like)
	} else {
		return l.repo.DeleteCommentLike(like.CommentId, like.UserID)
	}
}
