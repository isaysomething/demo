package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	PostStatusDraft     = 0
	PostStatusPublished = 1

	PostTypePost     = 1
	PostTypePage     = 2
	PostTypeRevision = 3
)

type Post struct {
	ID        uint64       `db:"id" json:"id"`
	UserID    uint64       `db:"user_id" json:"user_id"`
	Title     string       `db:"title" json:"title"`
	Content   string       `db:"content" json:"content"`
	Status    int          `db:"status" json:"status"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
}

func GetPostCount(db *sqlx.DB) (count int, err error) {
	err = db.Get(&count, "SELECT count(*) FROM posts")
	return
}

func GetPosts(db *sqlx.DB, page, limit int) (posts []Post, err error) {
	posts = []Post{}
	err = db.Select(&posts, "SELECT * FROM posts ORDER BY created_at LIMIT ? OFFSET ?", limit, (page-1)*limit)
	return
}

func GetPost(db *sqlx.DB, id uint64) (Post, error) {
	post := Post{}
	err := db.Get(&post, "SELECT * FROM posts WHERE id=?", id)
	return post, err
}
