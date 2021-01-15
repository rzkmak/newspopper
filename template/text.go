package template

import "newspopper/model"

func TextParser(data model.Article) ([]byte, error) {
	msg := "Update: " + data.Url + "\n" +
		data.Title + "\n" +
		"Open now: " + data.Link
	return []byte(msg), nil
}
