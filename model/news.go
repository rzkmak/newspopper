package model

type SiteUpdate struct {
	Name     string
	Articles []Article
}

type Article struct {
	Url   string
	Title string
	Link  string
}
