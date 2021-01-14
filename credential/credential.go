package credential

import (
	"errors"
	"fmt"
	"newspopper/loader"
)

type Credential interface {
	Get(t, alias string) (string, error)
}

type Storage struct {
	Creds map[string]map[string]string
}

func (s Storage) Get(t, alias string) (string, error) {
	if s.Creds == nil {
		return "", errors.New("uninitialized credential object")
	}
	if _, found := s.Creds[t]; !found {
		return "", errors.New("uninitialized credential for type:" + t + " alias:" + alias)
	}
	if _, found := s.Creds[t][alias]; !found {
		return "", errors.New("uninitialized credential for type:" + t + " alias:" + alias)
	}
	return s.Creds[t][alias], nil
}

func NewCredentialStorage(credential loader.Credential) (Storage, error) {
	var creds map[string]map[string]string = make(map[string]map[string]string, 10)
	for _, v := range credential {
		if _, found := v["type"]; !found {
			return Storage{}, errors.New(fmt.Sprintf("not found type for credential: %v", v))
		}
		t := v["type"]
		if _, found := v["alias"]; !found {
			return Storage{}, errors.New(fmt.Sprintf("not found type for credential: %v", t))
		}
		alias := v["alias"]
		if _, found := v["token"]; !found {
			return Storage{}, errors.New(fmt.Sprintf("not found type for credential: %v", alias))
		}

		if _, found := creds[t]; !found {
			creds[t] = make(map[string]string, 10)
		}
		creds[t][alias] = v["token"]
	}
	return Storage{Creds: creds}, nil
}
