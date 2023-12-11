package service

import (
	"fmt"
	"foruum/models"
	"foruum/repository"
	"log"
	"strconv"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (p *PostService) CreatePost(post models.Post) error {
	// categories := []string{"IT", "Sport", "Education", "News"}
	// for _, w := range categories {
	// if post.Category != w {
	// 	fmt.Println(post.Category)

	// 	return fmt.Errorf("Uncorrect choose category")
	// }

	return p.repo.CreatePost(post)
}

func (p *PostService) ShowAllPosts() ([]models.Post, error) {
	posts, err := p.repo.ShowAllPosts()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) GetPostByID(id string) (*models.Post, error) {
	_, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post, err := p.repo.GetPostByID(id)
	if err != nil {
		fmt.Printf("service: %s", err)
		return nil, err
	}
	return post, nil
}

func (p *PostService) ChangePost(newpost, oldPost models.Post, user models.User) error {
	if user.Username != oldPost.Author {
		return fmt.Errorf("Uncorrect change post author ")
	}
	if err := p.repo.ChangePost(newpost, oldPost.Id); err != nil {
		return fmt.Errorf("CHANGE :%w", err)
	}
	return nil
}

func (p *PostService) ShowMyPosts(userId int) ([]models.Post, error) {
	posts, err := p.repo.ShowMyPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) ShowMyCommentPosts(userId int) ([]models.Post, error) {
	posts, err := p.repo.ShowMyCommentPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) ShowMyLikedPosts(userId int) ([]models.Post, error) {
	posts, err := p.repo.ShowMyLikedPosts(userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) GetPostsByCategoty(category []string) ([]models.Post, error) {
	posts, err := p.repo.GetPostsByCategoty(category)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return posts, nil
}

func (p *PostService) DeletePost(userId int, post models.Post) error {
	if userId != post.AuthorId {
		return fmt.Errorf("you can't delete this post! because you aren't the author!")
	}
	if err := p.repo.DeletePost(post.Id); err != nil {
		return err
	}
	return nil
}
