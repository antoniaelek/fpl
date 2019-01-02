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
	updatePlayers(&scraperConfig)
	updateTeams(&scraperConfig)
	updateScores(&scraperConfig, gameweek)
}

// Scrape scores & update database
func updateScores(scraperConfig *scraper.Config, gameweek int) {
	scores := scraper.ScrapeLiveScores(scraperConfig, gameweek)
	err := scraper.RefreshScoresBucket(scraperConfig, scores)
	if err != nil {
		log.Fatalln("Error inserting scores: ", err)
	} else {
		log.Println("Parsed and inserted scores.")
	}
}

// Scrape teams & update database
func updateTeams(scraperConfig *scraper.Config) {
	teams, err := scraper.ScrapeTeams(scraperConfig)
	if err != nil {
		log.Fatalln("Error parsing teams: ", err)
	}
	err = scraper.RefreshTeamsBucket(scraperConfig, teams)
	if err != nil {
		log.Fatalln("Error inserting teams: ", err)
	} else {
		log.Println("Parsed and inserted teams.")
	}
}

// Scrape players & update database
func updatePlayers(scraperConfig *scraper.Config) {
	players, err := scraper.ScrapePlayers(scraperConfig)
	if err != nil {
		log.Fatalln("Error parsing players: ", err)
	}
	err = scraper.RefreshPlayersBucket(scraperConfig, players)
	if err != nil {
		log.Fatalln("Error inserting players: ", err)
	} else {
		log.Println("Parsed and inserted players.")
	}
}
