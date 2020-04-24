package post

type Form struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	State   int    `json:"state"`
}
