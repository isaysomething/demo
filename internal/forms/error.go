package forms

type Error struct {
}

func (err Error) Error() string {
	return ""
}
