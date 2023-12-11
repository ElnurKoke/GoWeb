package service

import (
	"foruum/models"
	"foruum/repository"
)

type Authorization interface {
	CreateUser(models.User) error
	GenerateToken(login, password string) (models.User, error)
	GetUserByToken(string) (models.User, error)
	DeleteToken(string) error
}

type Post interface {
	CreatePost(models.Post) error
	ShowAllPosts() ([]models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	GetPostsByCategoty([]string) ([]models.Post, error)
	ChangePost(newPost, oldPost models.Post, user models.User) error
	ShowMyPosts(userId int) ([]models.Post, error)
	ShowMyCommentPosts(userId int) ([]models.Post, error)
	ShowMyLikedPosts(userId int) ([]models.Post, error)
	DeletePost(int, models.Post) error
}
type Comment interface {
	CreateComment(models.Comment) error
	GetCommentByPostID(int) (*[]models.Comment, error)
}

type Like interface {
	SetPostLike(models.Like) error
	SetCommentLike(models.Like) error
}

type Dislike interface {
	SetPostDislike(models.Dislike) error
	SetCommentDislike(models.Dislike) error
}

type Service struct {
	Authorization
	Post
	Comment
	Like
	Dislike
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Post:          NewPostService(repos.Post),
		Comment:       NewCommentService(repos.Comment),
		Like:          NewLikeService(repos.Like),
		Dislike:       NewDislikeService(repos.Dislike),
	}
}
