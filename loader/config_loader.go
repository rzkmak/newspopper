package loader

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type RedisConf struct {
	Uri             string `yaml:"uri"`
	PersistDuration string `yaml:"persist_duration"`
}

type BackendConf struct {
	RedisConf `yaml:"redis"`
}

type Credential []map[string]string
type Output []map[string]string
type Listener []map[string]interface{}

type Config struct {
	Port       int         `yaml:"port"`
	Backend    BackendConf `yaml:"backend"`
	Credential Credential  `yaml:"credential"`
	Output     Output      `yaml:"output"`
	Listener   Listener    `yaml:"listener"`
}

func Load() (Config, error) {
	f, err := ioutil.ReadFile("./sites.yaml")
	if err != nil {
		return Config{}, err
	}

	var sourceSites Config
	if err := yaml.Unmarshal(f, &sourceSites); err != nil {
		return Config{}, err
	}

	return sourceSites, nil
}
