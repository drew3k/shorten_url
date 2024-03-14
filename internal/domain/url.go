package domain

type URL struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

type ShortenedURLList struct {
	URLs []URL
}
