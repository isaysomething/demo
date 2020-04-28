package rbac

import "time"

const (
	TypeRole = iota + 1
	TypePermission
)

type Item struct {
	ID        string    `db:"id" json:"id"`
	GroupID   string    `db:"group_id" json:"group_id"`
	Name      string    `db:"name" json:"name"`
	Type      int       `db:"item_type" json:"item_type"`
	Reserved  bool      `db:"reserved" json:"reserved"`
	Obj       string    `db:"obj" json:"obj"`
	Act       string    `db:"act" json:"act"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (i Item) IsRole() bool {
	return i.Type == TypeRole
}

func (i Item) IsPermission() bool {
	return i.Type == TypePermission
}

type Group struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Permissions []Item    `db:"permissions" json:"permissions"`
}
