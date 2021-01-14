package loader

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type SourceSite struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Selector struct {
		Main  string `yaml:"main"`
		Title string `yaml:"title"`
		Link  string `yaml:"link"`
	}
	OptionalHttpCode int `yaml:"optional_http_code"`
}

func Load() ([]SourceSite, error) {
	f, err := ioutil.ReadFile("./sites.yaml")
	if err != nil {
		return nil, err
	}

	var sourceSites []SourceSite
	if err := yaml.Unmarshal(f, &sourceSites); err != nil {
		return nil, err
	}

	return sourceSites, nil
}
