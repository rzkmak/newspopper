package main

import (
    "anipokev2/cli"
    "os"
)

func main() {
    c := cli.NewCli(os.Args)
    c.Run()
}
