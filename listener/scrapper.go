package listener

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"newspopper/backend"
	"newspopper/model"
	"newspopper/util"
	"regexp"
	"strings"
	"time"
)

type ScrapperSelector struct {
	Main  string
	Title string
	Link  string
}

type Scrapper struct {
	Url                string
	OptionalHttpStatus int
	Selector           ScrapperSelector
	Backend            backend.Backend
	Output             io.Writer
	Interval           time.Duration
}

func (s *Scrapper) String() string {
	return fmt.Sprintf("scrapper_job:%v", s.Url)
}

func (s *Scrapper) Spawn(ctx context.Context) error {
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

func (s Scrapper) Run() {
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

func (s Scrapper) validate() error {
	if s.Url == "" {
		return errors.New("scrapper url shouldn't be null")
	}

	if s.Output == nil {
		return errors.New("scrapper output shouldn't be null")
	}

	if s.Interval < time.Minute {
		return errors.New("scrapper minimum interval is one minute")
	}
	return nil
}

func (s *Scrapper) update() ([]model.Article, error) {
	res, err := http.Get(s.Url)
	defer res.Body.Close()

	if err != nil {
		log.Errorln("error while visiting page: ", s.Url)
		return nil, err
	}

	log.Println(s.OptionalHttpStatus)

	if res.StatusCode != 200 && res.StatusCode != s.OptionalHttpStatus {
		log.Errorln("web page error code : ", s.Url, res.StatusCode)
		return nil, err
	}

	updates := make([]model.Article, 0)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Errorln("failed to parse web page : ", s.Url)
		return nil, err
	}

	doc.Find(s.Selector.Main).Children().Each(func(i int, selection *goquery.Selection) {
		update := new(model.Article)
		update.Title = selection.Find(s.Selector.Title).Text()
		plain, err := selection.Find(s.Selector.Link).Html()
		if err != nil {
			log.Error("error while parsing html attribute", err)
		}
		re := regexp.MustCompile("href=\"(.*?)\"")
		link := re.FindStringSubmatch(plain)
		if len(link) > 0 {
			l := strings.ReplaceAll(link[0], "href=", "")
			l = strings.ReplaceAll(l, "\"", "")
			update.Link = l
		}
		updates = append(updates, *update)
	})

	return updates, nil
}
