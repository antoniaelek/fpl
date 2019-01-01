package main

import (
	"fpl/scraper"
	"log"
	"os"
	"strconv"

	"github.com/tkanos/gonfig"
)

func main() {
	gameweek, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("Please pass in a valid gameweek number as program parameter:", err)
		return
	}

	config := scraper.Config{}
	err = gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatalln("Error fetching config:", err)
		return
	}

	scores := scraper.Scrape(&config, gameweek)
	err = scraper.RefreshStore(&config, scores)
	if err != nil {
		log.Fatalln("Error scraping data:", err)
		return
	} else {
		log.Println("Parsed and inserted data.")
	}
}
