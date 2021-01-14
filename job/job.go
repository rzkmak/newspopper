package job

import (
	"fmt"
	"newspopper/bot"
	"newspopper/config"
	"newspopper/model"
	"newspopper/scrapper"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Job struct {
	S      *scrapper.Scrapper
	Bot    *bot.Telegram
	Config *config.Config
}

func NewJob(s *scrapper.Scrapper, bot *bot.Telegram, config *config.Config) *Job {
	return &Job{S: s, Bot: bot, Config: config}
}
func (j *Job) Execute() {
	log.Infoln("starting scheduled job at :", time.Now())
	ticker := time.NewTicker(j.Config.FetchInterval)
	for {
		select {
		case <-ticker.C:
			updates := j.S.Scrap()
			if len(updates) == 0 {
				continue
			}
			j.emit(updates)
		}
	}
}

func (j *Job) emit(fansubs []model.SiteUpdate) {
	rds := redis.NewClient(&redis.Options{
		Addr: j.Config.RedisUri,
		DB:   0,
	})
	defer rds.Close()

	for _, fs := range fansubs {
		log.Infoln("updating anime from fansub :", fs.Name)
		for _, a := range fs.Articles {
			title := fmt.Sprintf("%v:%v", fs.Name, ToSnakeCase(a.Title))
			_, err := rds.Get(title).Result()
			if err != nil {
				rds.Set(title, 1, j.Config.RedisTimeout)
				j.Bot.Broadcast(fs.Name, a.Title, a.Link)
			}
		}
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
