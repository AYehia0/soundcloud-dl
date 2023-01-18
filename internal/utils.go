package internal

import (
	"fmt"
	"os"
	"sync"

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

// get all the available qualities inside the track
// used to choose a track to download based on the target quality
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
	if len(qualities) < 1 {
		fmt.Printf("%s No qualities available to download!:\n", theme.Red("[-]"))
		os.Exit(0)
	}
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

func getPlaylistDownloadTracks(soundData *soundcloud.SoundData, clientId string) [][]soundcloud.DownloadTrack {

	var wg sync.WaitGroup
	listDownloadTracks := make([][]soundcloud.DownloadTrack, 0)

	playlistTracks := soundcloud.GetPlaylistTracks(soundData, clientId)

	fmt.Printf("%s Playlist contains : %s track\n", theme.Green("[+]"), theme.Magenta(len(playlistTracks)))
	if !promptYesNo() {
		fmt.Printf("%s Exiting...\n", theme.Red("[-]"))
		os.Exit(0)
	}
	for i, t := range playlistTracks {
		wg.Add(1)

		go func(t soundcloud.SoundData) {
			defer wg.Done()
			dlTrack := soundcloud.GetFormattedDL(&t, clientId)
			listDownloadTracks = append(listDownloadTracks, dlTrack)
		}(t)
		fmt.Printf("%s  %v - %s \n", theme.Green("[+]"), theme.Red(i+1), theme.Magenta(t.Title))
	}
	wg.Wait()
	return listDownloadTracks
}

// get a final track to be downloaded
// if bestQuality is false it will prompt the user to choose a quality
func getTrack(downloadTracks []soundcloud.DownloadTrack, bestQuality bool) soundcloud.DownloadTrack {

	// show available download options
	qualities := getQualities(downloadTracks)
	if !bestQuality {
		defaultQuality = chooseQuality(qualities)
	} else {
		defaultQuality = getHighestQuality(qualities)
	}

	return chooseTrackDownload(downloadTracks, defaultQuality)

}

func promptYesNo() bool {
	prompt := promptui.Select{
		Label: "Download [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%s Download aborted!\n", theme.Red("[-]"))
		os.Exit(0)
	}
	return result == "Yes"
}
