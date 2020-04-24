package loader

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Fansubs struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Selector struct {
		Main  string `yaml:"main"`
		Title string `yaml:"title"`
		Link  string `yaml:"link"`
	}
	OptionalHttpCode int `yaml:"optional_http_code"`
}

func Load() ([]Fansubs, error) {
	f, err := ioutil.ReadFile("./fansubs.yaml")
	if err != nil {
		return nil, err
	}

	var fansubs []Fansubs
	if err := yaml.Unmarshal(f, &fansubs); err != nil {
		return nil, err
	}

	return fansubs, nil
}
