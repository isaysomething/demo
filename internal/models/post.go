package models

import (
	"time"

	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	PostStateDraft     = 0
	PostStatePublished = 1

	PostTypePost     = 1
	PostTypePage     = 2
	PostTypeRevision = 3
)

type Post struct {
	ID        int64          `db:"id" json:"id"`
	UserID    int64          `db:"user_id" json:"user_id"`
	Title     string         `db:"title" json:"title"`
	Content   string         `db:"content" json:"content"`
	State     int            `db:"state" json:"state"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt sqlex.NullTime `db:"updated_at" json:"updated_at"`
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Content, validation.Required),
		validation.Field(&p.State, validation.In(PostStateDraft, PostStatePublished)),
	)
}

func (p *Post) Save(db *sqlex.DB) (err error) {
	if err = p.Validate(); err != nil {
		return err
	}
	res, err := db.Exec(
		"INSERT INTO posts(id, user_id, title, content, state, created_at, updated_at) VALUES(null, ?, ?, ?, ?, ?, ?)",
		p.UserID, p.Title, p.Content, p.State, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return
	}
	p.ID, err = res.LastInsertId()
	return
}

func (p *Post) Update(db *sqlex.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE posts SET title=?, content=?, state=?, updated_at=? WHERE id=?",
		p.Title, p.Content, p.State, time.Now(), p.ID,
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

func (p *Post) Delete(db *sqlex.DB) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", p.ID)
	return err
}

func GetPost(db *sqlex.DB, id int64) (*Post, error) {
	post := new(Post)
	err := db.Get(post, "SELECT * FROM posts WHERE id=?", id)
	return post, err
}
