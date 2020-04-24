package config

import (
	"github.com/subosito/gotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		log.Println("testing")
		return
	}

	if err := gotenv.Load(); err != nil {
		log.Println("loading config from os environment variable")
	}
}

func GetInt(key string) int {
	i, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Printf("getting error from parsing config with key: %v and value: %v\n", key, i)
		log.Fatalln(err)
	}
	return i
}

func GetString(key string) string {
	return os.Getenv(key)
}
