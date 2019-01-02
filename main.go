package main

import (
	"fpl/scraper"
	"log"
	"os"
	"strconv"

	"github.com/tkanos/gonfig"
)

// Config represents application configuration
type Config struct {
	GameweekOneID int
	DatabaseFile  string
	Language      string
	AudioFolder   string
	Teams         []string
}

func main() {
	// Get gameweek
	gameweek := 0
	if len(os.Args) < 2 {
		log.Fatalln("Please pass in a valid gameweek number as program parameter.")
		return
	}
	gameweek, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("Please pass in a valid gameweek number as program parameter. ", err)
		return
	}

	// Get config
	config := Config{}
	err = gonfig.GetConf("config.json", &config)
	if err != nil {
		log.Fatalln("Error fetching config:", err)
		return
	}

	// Scrape
	scraperConfig := scraper.Config{DatabaseFile: config.DatabaseFile, GameweekOneID: config.GameweekOneID}
	scraper := scraper.NewScraper(&scraperConfig)
	scraper.ScrapeToDbFull(gameweek)
}
