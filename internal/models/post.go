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
	ID        uint64       `db:"id"`
	UserID    uint64       `db:"user_id"`
	Title     string       `db:"title"`
	Content   string       `db:"content"`
	Status    int          `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func GetPosts(db *sqlx.DB, offset, limit int) (posts []Post, err error) {
	err = db.Select(posts, "SELECT * FROM posts LIMIT ? OFFSET ?", limit, offset)
	return
}

func GetPost(db *sqlx.DB, id uint64) (Post, error) {
	post := Post{}
	err := db.Get(&post, "SELECT * FROM posts WHERE id=?", id)
	return post, err
}
