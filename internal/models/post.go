package models

const (
	PostStatusDraft     = 0
	PostStatusPublished = 1

	PostVisibilityPrivate = 0
	PostVisibilityPublic  = 1

	PostTypePost     = 1
	PostTypePage     = 2
	PostTypeRevision = 3
)

type Post struct {
	Model
	UserID     uint64 `gorm:"index;not null"`
	Title      string
	Excerpt    string `gorm:"not null;default:''"`
	Status     int    `gorm:"not null;default:0"`
	Visibility int    `gorm:"not null;default:1"`
	Type       int    `gorm:"not null;"`

	User User
}
