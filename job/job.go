package job

import (
	"anipokev2/bot"
	"anipokev2/config"
	"anipokev2/model"
	"anipokev2/scrapper"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
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
			updates, err := j.S.Scrap()
			if err != nil {
				log.Errorln("failed job at ", time.Now())
				continue
			}
			j.emit(updates)
		}
	}
}

func (j *Job) emit(fansubs []model.Fansub) {
	rds := redis.NewClient(&redis.Options{
		Addr: j.Config.RedisUri,
		DB:   0,
	})
	defer rds.Close()

	for _, fs := range fansubs {
		log.Infoln("updating anime from fansub :", fs.Name)
		for _, a := range fs.Anime {
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
