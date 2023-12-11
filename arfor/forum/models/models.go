package models

type Like struct {
	Id        int
	PostID    int
	UserID    int
	CommentId int
	Active    bool
}

type Dislike struct {
	Id        int
	PostID    int
	UserID    int
	CommentId int
	Active    bool
}

type Comment struct {
	Id      int
	Author  string
	Date    string
	PostId  int
	Text    string
	Like    int
	Dislike int
	UserId  int
}

type Post struct {
	Id          int
	Title       string
	Description string
	Like        int
	Dislike     int
	AuthorId    int
	Author      string
	Date        string
	Category    string
	Path        string
}
