package models

import (
	"time"

	"github.com/clevergo/demo/pkg/sqlex"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	PostStateDeleted   = 0
	PostStateDraft     = 1
	PostStatePublished = 2

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

func (p *Post) Delete(db *sqlex.DB) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", p.ID)
	return err
}

func GetPost(db *sqlex.DB, id int64) (*Post, error) {
	post := new(Post)
	err := db.Get(post, "SELECT * FROM posts WHERE id=?", id)
	return post, err
}
