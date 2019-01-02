package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// Config represents scraper configuration.
type Config struct {
	GameweekOneID int
	DatabaseFile  string
}

// Score represents goal score match event.
type Score struct {
	Minute           string
	GoalPlayerName   string
	AssistPlayerName string
	HomeTeam         string
	AwayTeam         string
	HomeTeamGoals    int
	AwayTeamGoals    int
	Gameweek         int
}

// Team represents a team.
type Team struct {
	Name      string
	ShortName string
}

// Player represents a player.
type Player struct {
	Name       string
	FirstName  string
	SecondName string
	WebName    string
}

// ScoreDbEntry represent goal score match event as stored in database.
type ScoreDbEntry struct {
	Score     *Score
	Timestamp time.Time
	Processed bool
}

// ScrapeTeams scrapes teams.
func ScrapeTeams(config *Config) (teams []Team, err error) {
	result, err := getJSON("http://fantasy.premierleague.com/drf/bootstrap-static")
	if err != nil {
		return nil, err
	}

	teams = make([]Team, 0)
	switch elements := result["teams"].(type) {
	case []interface{}:
		for _, elem := range elements {
			switch team := elem.(type) {
			case map[string]interface{}:
				t := Team{
					Name:      team["name"].(string),
					ShortName: team["short_name"].(string),
				}
				teams = append(teams, t)
			}
		}
	}

	return teams, err
}

// ScrapePlayers scrapes players.
func ScrapePlayers(config *Config) (players []Player, err error) {
	result, err := getJSON("http://fantasy.premierleague.com/drf/bootstrap-static")
	if err != nil {
		return nil, err
	}

	players = make([]Player, 0)
	switch elements := result["elements"].(type) {
	case []interface{}:
		for _, element := range elements {
			switch player := element.(type) {
			case map[string]interface{}:
				pl := Player{
					Name:       player["first_name"].(string) + " " + player["second_name"].(string),
					FirstName:  player["first_name"].(string),
					SecondName: player["second_name"].(string),
					WebName:    player["web_name"].(string),
				}
				players = append(players, pl)
			}
		}
	}

	return players, err
}

// ScrapeLiveScores scrapes gameweek live scores.
// Method takes two parameters: pointer to application config and gameweek to scrape.
// It returns slice of scraped score events in gameweek.
func ScrapeLiveScores(config *Config, gameweek int) []Score {
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
				score.GoalPlayerName = player
				score.AssistPlayerName = assist
				result = append(result, score)
			}
		})
	})

	gameweekID := config.GameweekOneID + gameweek - 1
	c.Visit("https://www.premierleague.com/matchweek/" + strconv.Itoa(gameweekID) + "/blog")

	return result
}

func getJSON(url string) (result map[string]interface{}, err error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Test")

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Cannot fetch URL %q: %v", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unexpected http GET status: %s", resp.Status)
		return
	}

	// body, err := ioutil.ReadAll(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&result)
	// err = json.Unmarshal(body, &result)
	if err != nil {
		err = fmt.Errorf("Cannot decode JSON: %v", err)
	}

	return
}
