package cli

import (
    "anipokev2/config"
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

    log.Infoln("hello world")
}

