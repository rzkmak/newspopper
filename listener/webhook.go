package listener

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"newspopper/backend"
	"newspopper/model"
	"newspopper/template"
	"newspopper/util"
)

type Webhook struct {
	WebhookUrl string
	Parser     template.Parser
	Backend    backend.Backend
	Output     io.Writer
}

type WebhookRequestContract struct {
	Provider string `json:"provider"`
	Title    string `json:"title"`
	Link     string `json:"link"`
}

func (w Webhook) String() string {
	return fmt.Sprintf("webhook_job:")
}

func (w Webhook) Spawn(ctx context.Context) error {
	if err := w.validate(); err != nil {
		return err
	}
	http.HandleFunc(w.WebhookUrl, func(writer http.ResponseWriter, request *http.Request) {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		var req WebhookRequestContract
		if err := json.Unmarshal(body, &req); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
		}

		if err := w.validateHttpRequest(req); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		isUpdated, err := w.Backend.Get(fmt.Sprintf("%s:%s:%s", w, req.Provider, util.ToSnakeCase(req.Title)))
		if isUpdated == 0 || err != nil {
			if err := w.Backend.Set(fmt.Sprintf("%s:%s:%s", w, req.Provider, util.ToSnakeCase(req.Title))); err != nil {
				log.Infoln("webhook failed to set backend: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
			msg, err := w.Parser(model.Article{
				Url:   req.Provider,
				Title: req.Title,
				Link:  req.Link,
			})
			if err != nil {
				log.Infoln("webhook failed to parse update: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
			if _, err := w.Output.Write(msg); err != nil {
				log.Infoln("rss failed to send output: ", err)
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
		}
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("event received"))
	})
	return nil
}

func (w Webhook) validateHttpRequest(contract WebhookRequestContract) error {
	if contract.Title == "" {
		return errors.New("field(title) should be specified")
	}
	if contract.Link == "" {
		return errors.New("field(link) should be specified")
	}
	if contract.Provider == "" {
		return errors.New("field(provider) should be specified")
	}
	return nil
}

func (w Webhook) validate() error {
	if w.WebhookUrl == "" {
		return errors.New("url should be defined for webhook receiver")
	}
	return nil
}
