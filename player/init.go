package player

import "fmt"

// Config represents player configuration
type Config struct {
	Language    string
	AudioFolder string
	Teams       []string
}

// NewPlayer is constructor for Player type.
func NewPlayer(config *Config) *Player {
	player := Player{AudioFolder: config.AudioFolder, Language: config.Language}

	// Scores audio files
	for i := 0; i <= 10; i++ {
		player.PlayText(fmt.Sprint(i))
	}

	// Teams audio files
	for _, team := range config.Teams {
		player.PlayText(team)
	}

	return &player
}
