package internal

import (
	"fmt"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
)

var (
	defaultQuality = "low"
	soundData      = &soundcloud.SoundData{}
	SearchLimit    = 6
	offset         = 0
)

func Sc(args []string, downloadPath string, bestQuality bool, search bool) {

	url := ""
	if len(args) > 0 {
		url = args[0]
	}

	if url != "" && !initValidations(url) {
		return
	}

	clientId := soundcloud.GetClientId(url)

	if clientId == "" {
		fmt.Println("Something went wrong while getting the Client Id!")
		return
	}
	// --search-and-download
	if search {
		keyword := getUserSearch()
		apiUrl := soundcloud.GetSeachAPIUrl(keyword, SearchLimit, offset, clientId)
		searchResult := soundcloud.SearchTracksByKeyword(apiUrl, keyword, offset, clientId)

		// select one to download
		soundData = selectSearchUrl(searchResult)
	} else {

		apiUrl := soundcloud.GetTrackInfoAPIUrl(url, clientId)
		soundData = soundcloud.GetSoundMetaData(apiUrl, url, clientId)
		if soundData == nil {
			fmt.Printf("%s URL : %s doesn't return a valid track. Track is publicly accessed ?", theme.Red("[+]"), theme.Magenta(url))
			return
		}

		fmt.Printf("%s Track found. Title : %s - Duration : %s\n", theme.Green("[+]"), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))
	}

	list := soundcloud.GetFormattedDL(soundData, clientId)

	// show available download options
	qualities := getQualities(list)
	if !bestQuality {
		defaultQuality = chooseQuality(qualities)
	} else {
		defaultQuality = getHighestQuality(qualities)
	}

	track := chooseTrackDownload(list, defaultQuality)
	filePath := soundcloud.Download(track, downloadPath)

	// add tags
	err := soundcloud.AddMetadata(track, filePath)
	if err != nil {
		fmt.Printf("Error happend while adding tags to the track : %s\n", err)
	}
	fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(filePath))
}
