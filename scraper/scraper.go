package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gocolly/colly"
)

// Config represents application configuration
type Config struct {
	GameweekOneID int
	DatabaseFile  string
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
	Gameweek      int
}

// ScoreDbEntry represent goal score match event as stored in database.
type ScoreDbEntry struct {
	Score     *Score
	Timestamp time.Time
	Processed bool
}

// Scrape scrapes gameweek live scores
// Takes pointer to Config and gameweek to scrape
// Returns slice of scraped score events
func Scrape(config *Config, gameweek int) []Score {
	result := make([]Score, 0)

	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Fatalln("Something went wrong:", err)
	})

	// Fetch info for each goal event entry and add it to the result array
	c.OnHTML(".event.popUp.goal", func(el *colly.HTMLElement) {
		score := Score{}
		score.Gameweek = gameweek
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

// RefreshStore updates scores data in database store
// Takes reference to config and slice of score events
func RefreshStore(config *Config, scores []Score) error {
	// No data
	if scores == nil || len(scores) == 0 {
		return nil
	}

	// Open database, create it if it doesn't exist
	db, err := bolt.Open(config.DatabaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln("Something went wrong:", err)
	}
	defer db.Close()

	// Update data in database
	db.Update(func(tx *bolt.Tx) error {
		// Try to create bucket
		b, err := tx.CreateBucketIfNotExists([]byte(strconv.Itoa(scores[0].Gameweek)))
		if err != nil {
			return fmt.Errorf("Create bucket error: %s", err)
		}

		// Fill bucket
		for i := 0; i < len(scores); i++ {
			score := scores[i]

			// Calculate key
			key := score.Minute + "_" + score.HomeTeam + "_" + strconv.Itoa(score.HomeTeamGoals) + "_" + score.AwayTeam + "_" + strconv.Itoa(score.AwayTeamGoals)

			// Check if element with this key exists
			value := b.Get([]byte(key))
			if value != nil {
				// TODO check if something's changed in data
				continue
			}

			// Serialize value
			scoreDb := ScoreDbEntry{Score: &score, Timestamp: time.Now(), Processed: false}
			value, err = json.Marshal(scoreDb)
			if err != nil {
				return fmt.Errorf("Json serialization error: %s", err)
			}

			// Add the new key-value pair
			err = b.Put([]byte(key), []byte(value))
			if err != nil {
				return fmt.Errorf("Update bucket error: %s", err)
			}
		}
		return nil
	})

	// If any error occurred return it
	return err
}
