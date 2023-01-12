package internal

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/AYehia0/soundcloud-dl/pkg/client"
	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
	"github.com/PuerkitoBio/goquery"
)

var defaultQuality string = "low"

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

func Sc(args []string) {

	url := args[0]
	downloadPath := args[len(args)-1]

	if !initValidations(url) {
		return
	}

	statusCode, body, err := client.Get(url)

	if err != nil {
		log.Fatalf("An Error : %s happended while requesting : %s", err, url)
	}
	if statusCode != http.StatusOK {
		fmt.Println("URL doesn't exist : status not 200.")
		return
	}

	// convert the bytes array into something we can read, as goquery doesn't accept strings
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	soundData := soundcloud.GetSoundMetaData(doc)
	if soundData == nil {
		fmt.Printf("%s URL : %s doesn't return a valid track. Track is publicly accessed ?", theme.Red("[+]"), theme.Magenta(url))
		return
	}

	fmt.Printf("%s Track found. Title : %s - Duration : %s\n", theme.Green("[+]"), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))

	clientId := soundcloud.GetClientId(doc)
	list := soundcloud.GetFormattedDL(soundData.Transcodings, clientId)

	// show available download options
	printQualities(list)
	soundcloud.Download(chooseTrackDownload(list, defaultQuality), downloadPath)
	fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
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

func printQualities(tracks []soundcloud.DownloadTrack) {
	fmt.Printf("%s Available Qualities :\n", theme.Green("[+]"))
	for i, track := range tracks {
		// check the default one
		if track.Quality == defaultQuality {
			fmt.Printf("\t[%s]- %s %s\n", theme.Magenta(i+1), theme.Green(track.Quality), theme.Magenta("(selected)"))
		} else {
			fmt.Printf("\t[%s]- %s\n", theme.Magenta(i+1), theme.Green(track.Quality))
		}
	}
}
