package main

import (
	"log"
	"newspopper/cli"
	"os"
)

func main() {
	c := cli.NewCli(os.Args)
	if err := c.Run(); err != nil {
		log.Fatalln(err)
	}
}
