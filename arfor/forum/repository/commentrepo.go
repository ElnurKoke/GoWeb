package repository

import (
	"database/sql"
	"fmt"
	"foruum/models"
)

type Comment interface {
	CreateComment(models.Comment) error
	GetCommentByPostID(int) (*[]models.Comment, error)
}
type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (c *CommentRepo) CreateComment(comment models.Comment) error {
	query := `INSERT INTO comment (postId, userId, author, text, date) VALUES ($2, $3, $4, $5, $6)`
	_, err := c.db.Exec(query, comment.PostId, comment.UserId, comment.Author, comment.Text, comment.Date)
	if err != nil {
		fmt.Printf("repo: %s\n", err)
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (c *CommentRepo) GetCommentByPostID(postid int) (*[]models.Comment, error) {
	query := `select * from comment where postId=$1`
	rows, err := c.db.Query(query, postid)
	if err != nil {
		return nil, fmt.Errorf("get post comment: %w", err)
	}
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Author, &comment.Text, &comment.Like, &comment.Dislike, &comment.Date); err != nil {
			fmt.Printf("repo: %s\n", err)
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return &comments, nil
}
