package output

import (
	"errors"
	"fmt"
	"io"
	"newspopper/credential"
	"newspopper/loader"
)

type Output interface {
	Get(alias string) (io.Writer, error)
}

type Impl struct {
	Outputs map[string]io.Writer
}

func NewOutputs(creds credential.Credential, outs loader.Output) (Output, error) {
	var result map[string]io.Writer = map[string]io.Writer{}
	for idx, outs := range outs {
		if _, found := outs["alias"]; !found {
			return nil, errors.New(fmt.Sprintf("output doesn't have alias at index %v", idx))
		}

		if _, found := outs["type"]; !found {
			return nil, errors.New(fmt.Sprintf("output doesn't have type at index %v", idx))
		}

		t := outs["type"]
		switch t {
		case "telegram-channel":
			alias, found := outs["alias"]
			if !found {
				return nil, errors.New(fmt.Sprintf("output doesn't have alias at index %v", idx))
			}
			cred, found := outs["credential"]
			if !found {
				return nil, errors.New(fmt.Sprintf("output doesn't have credential for alias %v", alias))
			}
			token, err := creds.Get("telegram", cred)
			if err != nil {
				return nil, err
			}
			chatId, found := outs["channel"]
			if !found {
				return nil, errors.New(fmt.Sprintf("output doesn't have channel for alias %v", alias))
			}
			result[alias] = NewTelegramChannel(token, chatId)
		default:
			return nil, errors.New(fmt.Sprintf("output doesnt supported type: %v", outs["type"]))
		}
	}
	return Impl{Outputs: result}, nil
}

func (i Impl) Get(alias string) (io.Writer, error) {
	if _, found := i.Outputs[alias]; !found {
		return nil, errors.New("output didn't found for alias: " + alias)
	}
	return i.Outputs[alias], nil
}
