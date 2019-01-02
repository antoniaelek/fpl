package scraper

import "log"

// Config represents scraper configuration.
type Config struct {
	GameweekOneID int
	DatabaseFile  string
}

// Scraper represents scraper.
type Scraper struct {
	GameweekOneID int
	DatabaseFile  string
}

// NewScraper is constructor for Scraper type.
func NewScraper(config *Config) *Scraper {
	scraper := Scraper{DatabaseFile: config.DatabaseFile, GameweekOneID: config.GameweekOneID}
	scraper.updateTeams()
	scraper.updatePlayers()
	return &scraper
}

// Scrapes scores & updates scores bucket
func (scraper *Scraper) updateScores(gameweek int) {
	scores := scraper.ScrapeInMemory(gameweek)
	err := scraper.refreshScoresBucket(scores)
	if err != nil {
		log.Fatalln("Error inserting scores: ", err)
	} else {
		log.Println("Parsed and inserted scores.")
	}
}

// Scrape teams & update teams bucket
func (scraper *Scraper) updateTeams() {
	teams, err := scraper.scrapeTeams()
	if err != nil {
		log.Fatalln("Error parsing teams: ", err)
	}
	err = scraper.refreshTeamsBucket(teams)
	if err != nil {
		log.Fatalln("Error inserting teams: ", err)
	} else {
		log.Println("Parsed and inserted teams.")
	}
}

// Scrape players & update players bucket
func (scraper *Scraper) updatePlayers() {
	players, err := scraper.scrapePlayers()
	if err != nil {
		log.Fatalln("Error parsing players: ", err)
	}
	err = scraper.refreshPlayersBucket(players)
	if err != nil {
		log.Fatalln("Error inserting players: ", err)
	} else {
		log.Println("Parsed and inserted players.")
	}
}
