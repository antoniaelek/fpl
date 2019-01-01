package scraper

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Config represents application configuration
type Config struct {
	GameweekOneID int
}

// Score represents goal score match event.
type Score struct {
	Minute        string
	Goal          string
	Assist        string
	HomeTeam      string
	AwayTeam      string
	HomeTeamGoals int
	AwayTeamGoals int
}

// Scrape scrapes gameweek live scores
func Scrape(config *Config, gameweek int) []Score {
	result := make([]Score, 0)

	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Fatalln("Something went wrong:", err)
	})

	c.OnHTML(".event.popUp.goal", func(el *colly.HTMLElement) {
		score := Score{}
		el.ForEach(".row", func(i int, e *colly.HTMLElement) {
			if i == 0 {
				arr := strings.Split(e.Text, " ")
				if len(arr) >= 10 {
					score.HomeTeam = arr[1]
					score.HomeTeamGoals, _ = strconv.Atoi(arr[3])
					score.AwayTeamGoals, _ = strconv.Atoi(arr[7])
					score.AwayTeam = arr[9]
				}
			} else if i == 1 {
				min := e.ChildText(".min")
				player := e.ChildText(".player")
				if len(player) > 8 && player[len(player)-7:len(player)] == "( pen )" {
					player = player[0 : len(player)-8]
				}
				assist := e.ChildText(".assist")
				if len(assist) > 0 && assist[0:4] == "Ast." {
					assist = assist[5:len(assist)]
				}
				score.Minute = min
				score.Goal = player
				score.Assist = assist
				result = append(result, score)
			}
		})
	})

	gameweekID := config.GameweekOneID + gameweek - 1
	c.Visit("https://www.premierleague.com/matchweek/" + strconv.Itoa(gameweekID) + "/blog")

	return result
}
