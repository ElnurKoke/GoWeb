package repository

import (
	"database/sql"
	"fmt"
	"foruum/models"
	"strings"
)

var (
	rows *sql.Rows
	err  error
)

type Post interface {
	CreatePost(models.Post) error
	ShowAllPosts() ([]models.Post, error)
	GetPostByID(string) (*models.Post, error)
	GetPostsByCategoty([]string) ([]models.Post, error)
	ChangePost(models.Post, int) error
	ShowMyPosts(int) ([]models.Post, error)
	ShowMyCommentPosts(int) ([]models.Post, error)
	ShowMyLikedPosts(int) ([]models.Post, error)
	DeletePost(int) error
}

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (p *PostRepo) CreatePost(post models.Post) error {
	query := `INSERT INTO post (user_id, author, title, description, date,category) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = p.db.Exec(query, post.AuthorId, post.Author, post.Title, post.Description, post.Date, post.Category)
	if err != nil {

		return fmt.Errorf("create post: %w", err)
	}
	return nil
}

func (p *PostRepo) ShowAllPosts() ([]models.Post, error) {
	rows, err := p.db.Query(`select * from post`)
	if err != nil {
		return nil, fmt.Errorf("show all post: %w", err)
	}

	var posts []models.Post
	for rows.Next() {
		post := new(models.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("show all post, scan: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (p *PostRepo) GetPostByID(id string) (*models.Post, error) {
	rows, err := p.db.Query(`SELECT * FROM post WHERE id=?`, id)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	var post models.Post
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("get post by id: scan: %w", err)
		}
	}
	return &post, nil
}

func (p *PostRepo) ChangePost(post models.Post, oldPostId int) error {
	query := `UPDATE post SET title=$1,description=$2 where id=$3;`
	_, err := p.db.Exec(query, post.Title, post.Description, oldPostId)
	if err != nil {
		return fmt.Errorf("Uncorrect change post way: %w", err)

	}
	return nil

}

func (p *PostRepo) ShowMyPosts(userId int) ([]models.Post, error) {
	rows, err := p.db.Query(`select * from post where user_id=?`, userId)
	if err != nil {
		return nil, fmt.Errorf("show my post: %w", err)
	}

	var posts []models.Post
	for rows.Next() {
		post := new(models.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("show my post, scan: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (p *PostRepo) ShowMyCommentPosts(userId int) ([]models.Post, error) {
	rows, err := p.db.Query(`select * from post where id in (select postId from comment where userId=?)`, userId)
	if err != nil {
		return nil, fmt.Errorf("show my comment post: %w", err)
	}

	var posts []models.Post
	for rows.Next() {
		post := new(models.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("show my comment post, scan: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (p *PostRepo) ShowMyLikedPosts(userId int) ([]models.Post, error) {
	rows, err := p.db.Query(`select * from post where id in (select postId from like where userId=? and postId not null)`, userId)
	if err != nil {
		return nil, fmt.Errorf("show my liked post: %w", err)
	}

	var posts []models.Post
	for rows.Next() {
		post := new(models.Post)
		if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
			return nil, fmt.Errorf("show my liked post, scan: %w", err)
		}
		posts = append(posts, *post)
	}
	return posts, nil
}

func (p *PostRepo) GetPostsByCategoty(category []string) ([]models.Post, error) {
	fmt.Println(strings.Join(category, ", "))
	var posts []models.Post
	post := new(models.Post)
	for _, i := range category {
		rows, err := p.db.Query(`select * from post where category=?`, i)
		if err != nil {
			return nil, fmt.Errorf("show posts by category: %w", err)
		}
		for rows.Next() {
			if err = rows.Scan(&post.Id, &post.AuthorId, &post.Author, &post.Title, &post.Description, &post.Like, &post.Dislike, &post.Date, &post.Category); err != nil {
				return nil, fmt.Errorf("show posts by category scan: %w", err)
			}
			posts = append(posts, *post)
		}
	}
	return posts, nil
}

func (p *PostRepo) DeletePost(PostID int) error {
	query := `delete from post where  id = ?`
	_, err := p.db.Exec(query, PostID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}
