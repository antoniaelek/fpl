package main

import (
	"fmt"
	"fpl/scraper"
	"log"

	"github.com/tkanos/gonfig"
)

func main() {
	config := scraper.Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatalln("Something went wrong:", err)
		return
	}

	scores := scraper.Scrape(&config, 21)
	scraper.RefreshStore(&config, scores)
	fmt.Println(len(scores))
}
