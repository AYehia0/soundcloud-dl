package internal

import (
	"fmt"
	"os"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
	"github.com/manifoldco/promptui"
)

func initValidations(url string) bool {

	fmt.Printf("%s Validating the URL : %s\n", theme.Yellow("[+]"), theme.Magenta(url))

	// check if the url is valid
	if !soundcloud.IsValidUrl(url) {
		fmt.Printf("%s The Url : %s isn't a valid soundcloud URL\n", theme.Red("[+]"), theme.Magenta(url))
		return false
	}

	fmt.Printf("%s URL is Valid! \n", theme.Green("[+]"))
	fmt.Printf("%s Checking The Track on Soundcloud. \n", theme.Red("[+]"))

	return true

}

// TEMP: Just for now, return the quality
// the default quality is just mp3, highest is ogg
// if the quality doesn't exist return the first one
func chooseTrackDownload(tracks []soundcloud.DownloadTrack, target string) soundcloud.DownloadTrack {
	for _, track := range tracks {
		if track.Quality == target {
			return track
		}
	}
	return tracks[0]
}

func getQualities(tracks []soundcloud.DownloadTrack) []string {
	qualities := make([]string, 0)
	for _, track := range tracks {
		// check the default one
		qualities = append(qualities, track.Quality)
	}
	return qualities
}

// prompt select quality, only if --best not passed
func chooseQuality(qualities []string) string {
	fmt.Printf("%s Available Qualities :\n", theme.Green("[+]"))
	prompt := promptui.Select{
		Label: "Choose a quality:",
		Items: qualities,
	}
	_, q, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}
	return q
}

func getHighestQuality(qualities []string) string {
	allQualities := []string{"high", "medium", "low"}
	var in = func(a string, list []string) bool {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}

	for _, q := range allQualities {
		if in(q, qualities) {
			return q
		}
	}
	return ""
}

// select a url to download
func selectSearchUrl(searches *soundcloud.SearchResult) *soundcloud.SoundData {

	titles := make([]string, 0)

	for _, res := range searches.Sounds {
		titles = append(titles, res.Title)
	}
	prompt := promptui.Select{
		Label: "Choose a track to download :",
		Items: titles,
	}
	ind, _, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}

	return &searches.Sounds[ind]
}

// prompt the input search for a keyword
func getUserSearch() string {

	prompt := promptui.Prompt{
		Label:   "🔍Search: ",
		Default: "surah yasin",
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(0)
	}

	return result

}
