package oldmodels

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/clevergo/demo/pkg/sqlex"
)

const (
	PostMetaMarkdownContent = "markdown_content"
)

type PostMeta struct {
	PostID    int64          `db:"post_id"`
	Key       string         `db:"meta_key"`
	Value     string         `db:"meta_value"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt sqlex.NullTime `db:"updated_at" json:"updated_at"`
}

func GetPostMeta(db *sqlex.DB, id int64, key string) (*PostMeta, error) {
	sql, args, err := squirrel.Select("*").From("post_meta").Where(squirrel.Eq{
		"post_id":  id,
		"meta_key": key,
	}).ToSql()
	if err != nil {
		return nil, err
	}
	meta := new(PostMeta)
	err = db.Get(meta, sql, args...)
	return meta, err
}

func GetPostMetas(db *sqlex.DB, id int64) ([]PostMeta, error) {
	sql, args, err := squirrel.Select("*").From("post_meta").Where(squirrel.Eq{
		"post_id": id,
	}).ToSql()
	if err != nil {
		return nil, err
	}
	metas := []PostMeta{}
	err = db.Select(&metas, sql, args...)
	return metas, err
}
