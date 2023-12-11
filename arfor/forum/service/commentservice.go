package service

import (
	"foruum/models"
	"foruum/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (c *CommentService) CreateComment(commnet models.Comment) error {
	return c.repo.CreateComment(commnet)
}

func (c *CommentService) GetCommentByPostID(id int) (*[]models.Comment, error) {
	return c.repo.GetCommentByPostID(id)
}
