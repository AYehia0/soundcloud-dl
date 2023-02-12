package soundclouddl

import (
	"fmt"
	"os"
	"strings"
)

var Search bool
var DownloadPath string
var BestQuality bool
var TrackQuality Quality

type Quality string

const (
	Mp3Quality Quality = "mp3"
	OggQuality Quality = "ogg"
)

func (q *Quality) String() string {
	return string(*q)
}

func (q *Quality) Set(value string) error {
	quality := Quality(strings.ToLower(value))
	if quality != Mp3Quality && quality != OggQuality {
		return fmt.Errorf("Invalid quality value: %s (valid values are 'mp3' and 'ogg')", value)
	}
	*q = quality
	return nil
}

func (e *Quality) Type() string {
	return "quality"
}

// define flags and handle configuration
func InitConfigVars() {
	tmpDLdir, _ := os.Getwd()
	rootCmd.PersistentFlags().BoolVarP(&Search, "search-and-download", "s", false, "Search for tracks by title and prompt one for download ")
	rootCmd.PersistentFlags().StringVarP(&DownloadPath, "download-path", "p", tmpDLdir, "The download path where tracks are stored.")
	rootCmd.PersistentFlags().BoolVarP(&BestQuality, "best", "b", false, "Download with the best available quality.")
	rootCmd.Flags().VarP(&TrackQuality, "quality", "q", "Spcifiy a download quality (MP3/OGG).")
}
