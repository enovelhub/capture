package book

type Book struct {
	Author   string    `json:"author"`
	Name     string    `json:"title"`
	Chapters []Chapter `json:"content"`
}

type Chapter struct {
	Title   string `json:"title"`
	Content []string `json:"content"`
}
