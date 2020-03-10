package models

type Comment struct {
	Model
	PostID  uint64
	Content string `gorm:"text;not null"`
}
