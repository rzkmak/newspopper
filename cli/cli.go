package cli

import (
	"fmt"
	"newspopper/bot"
	"newspopper/config"
	"newspopper/job"
	"newspopper/loader"
	"newspopper/scrapper"
	"os"
	"time"

	tb "github.com/demget/telebot"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
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

	if len(fansubs) == 0 {
		log.Fatalln("no fansubs detected, please fill in the sites.yaml")
	}

	s := scrapper.NewScrapper(fansubs)

	if len(c.Args) > 1 && c.Args[1] == "simulate" {
		log.Infoln("starting simulation mode")
		fs := s.Scrap()
		if len(fs) == 0 {
			log.Fatalln("simulator failed: no value returned")
		}
		for _, v := range fs {
			fmt.Println("getting result from :", v.Name)
			for k, a := range v.Articles {
				fmt.Println(k)
				fmt.Println("title: ", a.Title)
				fmt.Println("link: ", a.Link)
			}
		}
		return
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
