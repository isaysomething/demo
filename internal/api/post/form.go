package post

import (
	"database/sql"
	"time"

	"github.com/clevergo/demo/internal/models"
	"github.com/clevergo/demo/pkg/sqlex"
)

type Form struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

func (f *Form) Create(db *sqlex.DB) (p *models.Post, err error) {
	now := time.Now()
	p = &models.Post{
		Title:     f.Title,
		Content:   f.Content,
		Status:    f.Status,
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}
	err = p.Save(db)
	return
}

func (f *Form) Update(db *sqlex.DB, p *models.Post) (err error) {
	now := time.Now()
	p.Title = f.Title
	p.Content = f.Content
	p.Status = f.Status
	p.UpdatedAt = sql.NullTime{Time: now, Valid: true}
	err = p.Update(db)
	return
}
