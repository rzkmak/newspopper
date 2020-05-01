package scrapper

import (
	"anipokev2/loader"
	"anipokev2/model"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Scrapper struct {
	Fs []loader.Fansubs
}

func NewScrapper(fs []loader.Fansubs) *Scrapper {
	return &Scrapper{Fs: fs}
}

func (s *Scrapper) Scrap() []model.Fansub {
	fansubs := make([]model.Fansub, 0)
	for _, fs := range s.Fs {
		anime, err := s.update(fs)
		if err != nil {
			log.Errorln("error while scrapping anime: ", err)
			log.Errorln("failed job at: ", time.Now())
			continue
		}
		fansub := new(model.Fansub)
		fansub.Name = fs.Name
		fansub.Anime = anime
		fansubs = append(fansubs, *fansub)
	}

	return fansubs
}

func (s *Scrapper) update(fs loader.Fansubs) ([]model.Anime, error) {
	res, err := http.Get(fs.Url)
	defer res.Body.Close()

	if err != nil {
		log.Errorln("error while visiting anime page: ", fs.Url, fs.Name)
		return nil, err
	}

	log.Println(fs.OptionalHttpCode)

	if res.StatusCode != 200 && res.StatusCode != fs.OptionalHttpCode {
		log.Errorln("web page error code : ", fs.Url, fs.Name, res.StatusCode)
		return nil, err
	}

	updates := make([]model.Anime, 0)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Errorln("failed to parse web page : ", fs.Url, fs.Name)
		return nil, err
	}

	doc.Find(fs.Selector.Main).Children().Each(func(i int, selection *goquery.Selection) {
		update := new(model.Anime)
		update.Title = selection.Find(fs.Selector.Title).Text()
		plain, err := selection.Find(fs.Selector.Link).Html()
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