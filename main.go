package main

import (
	"newspopper/cli"
	"os"
)

func main() {
	c := cli.NewCli(os.Args)
	c.Run()
}
