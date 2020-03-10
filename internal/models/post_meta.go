package models

const (
	PostMetaContent = "content"
)

type PostMeta struct {
	PostID    int    `gorm:"not null"`
	MetaKey   string `gorm:"varchar(64);not null"`
	MetaValue string `gorm:"not null"`
}
