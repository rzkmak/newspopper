package cli

import (
	"anipokev2/bot"
	"anipokev2/config"
	"anipokev2/job"
	"anipokev2/loader"
	"anipokev2/scrapper"
	"fmt"
	tb "github.com/demget/telebot"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type Cli struct {
	*config.Config
	Args []string
}

func NewCli(args []string) *Cli {
	return &Cli{Args: args}
}

func (c *Cli) Run() {
	log.SetLevel(log.DebugLevel)
	log.StandardLogger()
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	fansubs, err := loader.Load()
	if err != nil {
		log.Fatalln(err)
	}
	s := scrapper.NewScrapper(fansubs)

	if c.Args[1] == "simulate" {
		log.Infoln("starting simulation mode")
		fs, err := s.Scrap()
		if err != nil {
			log.Fatalln("simulator failed: ", err)
		}
		for _, v := range fs {
			fmt.Println("getting result from :", v.Name)
			for k, a := range v.Anime {
				fmt.Println(k)
				fmt.Println("title: ", a.Title)
				fmt.Println("link: ", a.Link)
			}
		}
		return
	}

	if len(fansubs) == 0 {
		log.Fatalln("no fansubs detected, please fill in the fansubs.yaml")
	}

	cfg := config.NewConfig()

	rds := redis.NewClient(&redis.Options{
		Addr: cfg.RedisUri,
		DB:   0,
	})
	defer rds.Close()
	if rds.Ping().Err() != nil {
		log.Fatalln(err)
	}

	p := &tb.LongPoller{Timeout: 15 * time.Second}

	t, err := tb.NewBot(tb.Settings{
		Token:  cfg.TelegramToken,
		Poller: p,
	})
	b := bot.NewTelegram(t, cfg, rds)
	go b.Run()

	scheduled := job.NewJob(s, b, cfg)
	scheduled.Execute()
}
