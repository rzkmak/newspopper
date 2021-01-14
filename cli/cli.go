package cli

import (
	"context"
	log "github.com/sirupsen/logrus"
	"newspopper/backend"
	"newspopper/credential"
	"newspopper/listener"
	"newspopper/loader"
	"newspopper/output"
	"os"
	"os/signal"
	"syscall"
)

type Cli struct {
	Args []string
}

func NewCli(args []string) *Cli {
	return &Cli{Args: args}
}

func (c *Cli) Run() error {
	config, err := loader.Load()
	if err != nil {
		return err
	}

	store, err := backend.NewBackend(config.Backend)
	if err != nil {
		return err
	}

	creds, err := credential.NewCredentialStorage(config.Credential)
	if err != nil {
		return err
	}

	outputs, err := output.NewOutputs(creds, config.Output)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	listeners, err := listener.NewListeners(config.Listener, store, outputs)
	if err != nil {
		return err
	}

	listeners.Initiate(ctx)
	waitForShutdown(ctx)
	return nil
}

func waitForShutdown(ctx context.Context) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	ctx.Done()
	log.Warn("Api server shutting complete")
}
