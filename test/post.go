package test

// Post ...
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Likes   int    `json:"likes"`
}
