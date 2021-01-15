package template

import (
	"encoding/json"
	"newspopper/model"
)

func JsonParser(data model.Article) ([]byte, error) {
	return json.Marshal(data)
}
