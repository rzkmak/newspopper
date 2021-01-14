package listener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"newspopper/backend"
	"newspopper/loader"
	"newspopper/output"
	"time"
)

type Listener interface {
	Spawn(ctx context.Context) error
}

type Listeners interface {
	Initiate(ctx context.Context)
}

type Impl struct {
	Listeners []Listener
}

func (i Impl) Initiate(ctx context.Context) {
	for _, listener := range i.Listeners {
		go func(l Listener) {
			if err := l.Spawn(ctx); err != nil {
				log.Fatalln(err)
			}
		}(listener)
	}
}

func NewListeners(listeners loader.Listener, backend backend.Backend, output output.Output) (Listeners, error) {
	result := []Listener{}
	for idx, listener := range listeners {
		if _, found := listener["type"]; !found {
			return nil, errors.New(fmt.Sprintf("listener doesn't have type at index %v", idx))
		}

		if _, found := listener["target"]; !found {
			return nil, errors.New(fmt.Sprintf("listener doesn't have target with spec %v", listener))
		}
		if _, found := listener["url"]; !found {
			return nil, errors.New(fmt.Sprintf("listener doesn't have url with spec %v", listener))
		}
		if _, found := listener["interval"]; !found {
			return nil, errors.New(fmt.Sprintf("listener doesn't have interval with spec %v", listener))
		}

		url := listener["url"].(string)
		target := listener["target"].(string)
		interval := listener["interval"].(string)
		t := listener["type"]
		switch t {
		case "scrapper":
			optionalHttpCode := 200
			if httpCode, found := listener["optional_http_code"]; found {
				optionalHttpCode = httpCode.(int)
			}
			if _, found := listener["selector"]; !found {
				return nil, errors.New(fmt.Sprintf("listener doesn't have selector with spec %v", listener))
			}
			selector := make(map[string]string)

			for key, value := range listener["selector"].(map[interface{}]interface{}) {
				strKey := fmt.Sprintf("%v", key)
				strValue := fmt.Sprintf("%v", value)

				selector[strKey] = strValue
			}
			if _, found := selector["main"]; !found {
				return nil, errors.New(fmt.Sprintf("listener doesn't have selector main with spec %v", listener))
			}
			if _, found := selector["title"]; !found {
				return nil, errors.New(fmt.Sprintf("listener doesn't have selector title with spec %v", listener))
			}
			if _, found := selector["link"]; !found {
				return nil, errors.New(fmt.Sprintf("listener doesn't have selector link with spec %v", listener))
			}

			out, err := output.Get(target)
			if err != nil {
				return nil, err
			}

			intervalDur, err := time.ParseDuration(interval)
			if err != nil {
				return nil, err
			}

			result = append(result, &Scrapper{
				Url:                url,
				OptionalHttpStatus: optionalHttpCode,
				Selector: ScrapperSelector{
					Main:  selector["main"],
					Title: selector["title"],
				},
				Backend:  backend,
				Output:   out,
				Interval: intervalDur,
			})

		default:
			return nil, errors.New(fmt.Sprintf("listener type %v unavailable with spec %v", t, listener))
		}
	}
	return Impl{Listeners: result}, nil
}
