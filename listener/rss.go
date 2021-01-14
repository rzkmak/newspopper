package listener

import (
	"context"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"io"
	"newspopper/backend"
	"newspopper/model"
	"newspopper/util"
	"time"
)

type Rss struct {
	Url      string
	Backend  backend.Backend
	Output   io.Writer
	Interval time.Duration
}

func (s *Rss) String() string {
	return fmt.Sprintf("rss_job:%v", s.Url)
}

func (s *Rss) Spawn(ctx context.Context) error {
	if err := s.validate(); err != nil {
		return err
	}
	log.Infoln("starting scrapper", s)

	tick := time.NewTicker(s.Interval)
	go func() {
		for {
			select {
			case t := <-tick.C:
				log.Infoln("running scrapper %s", s, t)
				s.Run()
			case <-ctx.Done():
				log.Infoln("shutdown scrapper %s", s)
				return
			}
		}
	}()
	return nil
}

func (s Rss) validate() error {
	if s.Url == "" {
		return errors.New("rss url shouldn't be null")
	}

	if s.Output == nil {
		return errors.New("rss output shouldn't be null")
	}

	if s.Interval < time.Minute {
		return errors.New("rss minimum interval is one minute")
	}
	return nil
}

func (s Rss) Run() {
	latest, err := s.update()
	if err != nil {
		log.Infoln("scrapper fetch error: ", err)
	}

	for _, l := range latest {
		isUpdated, err := s.Backend.Get(fmt.Sprintf("%s:%s", s, util.ToSnakeCase(l.Title)))
		if isUpdated == 0 || err != nil {
			if err := s.Backend.Set(fmt.Sprintf("%s:%s", s, util.ToSnakeCase(l.Title))); err != nil {
				log.Infoln("scrapper failed to set backend: ", err)
				continue
			}
			msg := "Update: " + s.Url + "\n" +
				l.Title + "\n" +
				"Open now: " + l.Link
			if _, err := s.Output.Write([]byte(msg)); err != nil {
				log.Infoln("scrapper failed to send output: ", err)
			}
		}
	}

}

func (s *Rss) update() ([]model.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(s.Url, ctx)
	if err != nil {
		log.Errorln("error while visiting page: ", s.Url)
		return nil, err
	}

	updates := make([]model.Article, 0)
	for _, item := range feed.Items {
		updates = append(updates, model.Article{
			Title: item.Title,
			Link:  item.Link,
		})
	}
	return updates, nil
}
