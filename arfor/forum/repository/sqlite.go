package repository

import (
	"database/sql"
	"log"
)

const (
	userBD = `CREATE TABLE IF NOT EXISTS user(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			password TEXT ,
			login TEXT UNIQUE,
			username TEXT,
			token TEXT DEFAULT NULL,
			tokenduration DATETIME DEFAULT NULL
			);`
	PostBD = `CREATE TABLE IF NOT EXISTS post(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			author TEXT,
			title TEXT,
			description TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			category TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id)
			);`
	CommentBD = `CREATE TABLE IF NOT EXISTS comment(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			postId INTEGER,
			userId INTEGER,
			author TEXT,
			text TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			FOREIGN KEY(postId) REFERENCES post(id)
			);`
	LikeBD = `CREATE TABLE IF NOT EXISTS like(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id)
		);`
	DislikeDB = `CREATE TABLE IF NOT EXISTS dislike(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id)
		);`
)

type ConfigDB struct {
	Path   string
	Driver string
	Name   string
}

func NewConfigDB() *ConfigDB {
	return &ConfigDB{
		Driver: "sqlite3",
		Name:   "forum.db",
		Path:   "repository",
	}
}

func InitDB(c *ConfigDB) (*sql.DB, error) {
	db, err := sql.Open(c.Driver, c.Name)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db sql.DB) error {
	for _, table := range []string{userBD, PostBD, CommentBD, LikeBD, DislikeDB} {
		_, err := db.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
