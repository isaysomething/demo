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
	ID        int64        `db:"id" json:"id"`
	UserID    uint64       `db:"user_id" json:"user_id"`
	Title     string       `db:"title" json:"title"`
	Content   string       `db:"content" json:"content"`
	Status    int          `db:"status" json:"status"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
}

func (p *Post) Save(db *sqlx.DB) error {
	res, err := db.Exec(
		"INSERT INTO posts(id, user_id, title, content, status, created_at, updated_at) VALUES(null, ?, ?, ?, ?, ?, ?)",
		p.UserID, p.Title, p.Content, p.Status, p.CreatedAt, time.Now(),
	)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) Update(db *sqlx.DB) error {
	res, err := db.Exec(
		"UPDATE posts SET title=?, content=?, status=?, updated_at=? WHERE id=?",
		p.Title, p.Content, p.Status, time.Now(), p.ID,
	)
	if err != nil {
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetPostsCount(db *sqlx.DB) (count int, err error) {
	err = db.Get(&count, "SELECT count(*) FROM posts")
	return
}

func GetPosts(db *sqlx.DB, limit, offset int) (posts []Post, err error) {
	posts = []Post{}
	err = db.Select(&posts, "SELECT * FROM posts ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	return
}

func GetPost(db *sqlx.DB, id int64) (*Post, error) {
	post := new(Post)
	err := db.Get(post, "SELECT * FROM posts WHERE id=?", id)
	return post, err
}
