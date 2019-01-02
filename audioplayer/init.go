package audioplayer

import "fmt"

// Config represents player configuration
type Config struct {
	Language    string
	AudioFolder string
	Teams       []string
}

// NewAudioPlayer is constructor for AudioPlayer type.
func NewAudioPlayer(config *Config) *AudioPlayer {
	player := AudioPlayer{AudioFolder: config.AudioFolder, Language: config.Language}

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
