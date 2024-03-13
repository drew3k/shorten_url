package domain

type URL struct {
	ID        int    `json:"id"`
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

type ShortenedURLList struct {
	URLs []URL
}
