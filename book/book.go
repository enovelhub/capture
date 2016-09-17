package book

type Book struct {
	Author string `json:"author"`
	Name string `json:"name"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	Title string `json:"title"`
	Content string `json:"content"`
}
