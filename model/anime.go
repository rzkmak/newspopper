package model

type Fansub struct {
	Name  string
	Anime []Anime
}

type Anime struct {
	Title string
	Link  string
}
