package main

import (
	"log"
	"os"

	"github.com/upframe/upframy"
)

func main() {
	token := os.Getenv("SLACK_BOT_TOKEN")

	if token == "" {
		panic("SLACK_BOT_TOKEN variable not defined")
	}

	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	upframy.Start(token, logger)
}
