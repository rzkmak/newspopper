package cli

import (
	"anipokev2/config"
	"anipokev2/loader"
	"anipokev2/scrapper"
	log "github.com/sirupsen/logrus"
	"os"
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
		log.Fatalln("no fansubs detected, please fill in the fansubs.yaml")
	}

	s := scrapper.NewScrapper(fansubs)
	result, err := s.Scrap()
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln(result)
}
