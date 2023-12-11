package service

import (
	"foruum/models"
	"foruum/repository"
)

type DislikeService struct {
	repo repository.Dislike
}

func NewDislikeService(repo repository.Dislike) *DislikeService {
	return &DislikeService{
		repo: repo,
	}
}

func (d *DislikeService) SetPostDislike(dislike models.Dislike) error {
	if d.repo.CheckPostLike(dislike.PostID, dislike.UserID) == nil {
		d.repo.DeletePostLike(dislike.PostID, dislike.UserID)
	}
	if d.repo.CheckPostDislike(dislike.PostID, dislike.UserID) != nil {
		return d.repo.SetPostDislike(dislike)
	} else {
		return d.repo.DeletePostDislike(dislike.PostID, dislike.UserID)
	}
}

func (d *DislikeService) SetCommentDislike(dislike models.Dislike) error {
	if d.repo.CheckCommentLike(dislike.CommentId, dislike.UserID) == nil {
		d.repo.DeleteCommentLike(dislike.CommentId, dislike.UserID)
	}
	if d.repo.CheckCommentDislike(dislike.CommentId, dislike.UserID) != nil {
		return d.repo.SetCommentDislike(dislike)
	} else {
		return d.repo.DeleteCommentDislike(dislike.CommentId, dislike.UserID)
	}
}
