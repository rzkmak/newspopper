package model

type SiteUpdate struct {
	Name     string
	Articles []Article
}

type Article struct {
	Title string
	Link  string
}
