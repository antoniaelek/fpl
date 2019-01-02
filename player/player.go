package player

import (
	"os/exec"

	htgotts "github.com/hegedustibor/htgo-tts"
)

// Player represents audio player.
type Player struct {
	Language    string
	AudioFolder string
}

// PlayText creates audio from the given text and plays it.
func (player *Player) PlayText(text string) {
	speech := htgotts.Speech{Folder: player.AudioFolder, Language: player.Language}
	speech.Speak(text)
}

// PlayFile plays audio file
func (player *Player) PlayFile(fileName string) error {
	mplayer := exec.Command("mplayer", player.AudioFolder+"/"+fileName+".mp3")
	return mplayer.Run()
}

// CreateAudioFiles creates audio files for given array of strings.
func CreateAudioFiles(player *Player, strings []string) {
	speech := htgotts.Speech{Folder: player.AudioFolder, Language: player.Language}
	for _, str := range strings {
		speech.Speak(str)
	}
}
