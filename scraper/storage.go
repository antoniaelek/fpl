package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// ScoreDbEntry represent goal score match event as stored in database.
type ScoreDbEntry struct {
	Score     *Score
	Timestamp time.Time
	Processed bool
}

// ScrapeToDb scrapes scores & updates database
func (scraper *Scraper) ScrapeToDb(gameweek int) {
	scraper.updateScores(gameweek)
}

// ScrapeToDbFull scrapes scores & updates database
func (scraper *Scraper) ScrapeToDbFull(gameweek int) {
	scraper.updateTeams()
	scraper.updatePlayers()
	scraper.updateScores(gameweek)
}

// Updates gameweek scores bucket.
func (scraper *Scraper) refreshScoresBucket(scores []Score) error {
	if scores == nil || len(scores) == 0 {
		return nil
	}

	interfaceSlice := make([]interface{}, len(scores))
	for i, el := range scores {
		interfaceSlice[i] = el
	}

	return scraper.refreshBucket(strconv.Itoa(scores[0].Gameweek), interfaceSlice, scoresBucketKey, scoresBuckeValue, scoresBucketValueFresh)
}

// Updates teams bucket.
func (scraper *Scraper) refreshTeamsBucket(teams []Team) error {
	if teams == nil || len(teams) == 0 {
		return nil
	}

	interfaceSlice := make([]interface{}, len(teams))
	for i, el := range teams {
		interfaceSlice[i] = el
	}

	return scraper.refreshBucket("teams", interfaceSlice, teamsBucketKey, teamsBucketValue, teamsBucketValueFresh)
}

// Updates players bucket.
func (scraper *Scraper) refreshPlayersBucket(players []Player) error {
	if players == nil || len(players) == 0 {
		return nil
	}

	interfaceSlice := make([]interface{}, len(players))
	for i, el := range players {
		interfaceSlice[i] = el
	}

	return scraper.refreshBucket("players", interfaceSlice, playersBucketKey, playersBucketValue, playersBucketValueFresh)
}

func scoresBucketKey(score interface{}) (result []byte, err error) {
	switch s := score.(type) {
	case Score:
		return []byte(s.Minute + "_" + s.HomeTeam + "_" + strconv.Itoa(s.HomeTeamGoals) + "_" + s.AwayTeam + "_" + strconv.Itoa(s.AwayTeamGoals)), nil
	}
	err = fmt.Errorf("Argument must be of type Score")
	return
}

func scoresBuckeValue(score interface{}) (result []byte, err error) {
	switch s := score.(type) {
	case Score:
		entry := ScoreDbEntry{
			Score:     &s,
			Timestamp: time.Now(),
			Processed: false,
		}
		result, err = json.Marshal(entry)
		return
	}
	err = fmt.Errorf("Argument must be of type Score")
	return
}

func scoresBucketValueFresh(value []byte, score interface{}) (isFresh bool, err error) {
	isFresh = false
	if value == nil {
		return
	}

	switch s := score.(type) {
	case Score:
		var data ScoreDbEntry
		err = json.Unmarshal(value, data)
		if err != nil {
			return
		}

		if data.Score.GoalPlayerName == s.GoalPlayerName && s.AssistPlayerName == data.Score.AssistPlayerName {
			isFresh = true
		}

		return
	}
	err = fmt.Errorf("Argument must be of type Score")
	return
}

func teamsBucketKey(team interface{}) (result []byte, err error) {
	switch t := team.(type) {
	case Team:
		return []byte(t.ShortName), nil
	}
	err = fmt.Errorf("Argument must be of type Team")
	return
}

func teamsBucketValue(team interface{}) (result []byte, err error) {
	switch t := team.(type) {
	case Team:
		return json.Marshal(t)
	}
	err = fmt.Errorf("Argument must be of type Team")
	return
}

func teamsBucketValueFresh(value []byte, team interface{}) (isFresh bool, err error) {
	isFresh = false
	if value == nil {
		return
	}

	switch team.(type) {
	case Team:
		return true, nil
	}
	err = fmt.Errorf("Argument must be of type Team")
	return
}

func playersBucketKey(player interface{}) (result []byte, err error) {
	switch p := player.(type) {
	case Player:
		return []byte(p.Name), nil
	}
	err = fmt.Errorf("Argument must be of type Player")
	return
}

func playersBucketValue(player interface{}) (result []byte, err error) {
	switch p := player.(type) {
	case Player:
		return json.Marshal(p)
	}
	err = fmt.Errorf("Argument must be of type Player")
	return
}

func playersBucketValueFresh(value []byte, player interface{}) (isFresh bool, err error) {
	isFresh = false
	if value == nil {
		return
	}

	switch player.(type) {
	case Player:
		return true, nil
	}
	err = fmt.Errorf("Argument must be of type Player")
	return
}

func (scraper *Scraper) refreshBucket(bucketName string, data []interface{},
	keySelector func(e interface{}) (result []byte, err error),
	valueSelector func(e interface{}) (result []byte, err error),
	valueFreshCheck func(value []byte, freshValue interface{}) (isFresh bool, err error)) error {
	// No data
	if data == nil || len(data) == 0 {
		return nil
	}

	// Open database, create it if it doesn't exist
	db, err := bolt.Open(scraper.DatabaseFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln("Error opening connection to database:", err)
		return err
	}
	defer db.Close()

	// Update data in database
	db.Update(func(tx *bolt.Tx) error {
		// Try to create bucket
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("Create bucket error: %s", err)
		}

		// Fill bucket
		for _, element := range data {
			// Calculate key
			key, err := keySelector(element)
			if err != nil {
				log.Println(fmt.Errorf("Key calculation error: %s", err))
				continue
			}

			// Check if element with this key exists and if it is fresh
			value := b.Get([]byte(key))
			isFresh, err := valueFreshCheck(value, element)
			if err != nil && isFresh {
				continue
			}

			// Create value
			value, err = valueSelector(element)
			if err != nil {
				log.Println(fmt.Errorf("Value calculation error: %s", err))
				continue
			}

			// Add the new key-value pair
			err = b.Put([]byte(key), []byte(value))
			if err != nil {
				log.Println(fmt.Errorf("Update bucket error: %s", err))
			}
		}
		return nil
	})

	// If any error occurred return it
	return err
}
