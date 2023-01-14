package internal

import (
	"fmt"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
)

var defaultQuality string = "low"

func Sc(args []string, bestQuality bool) {

	url := args[0]
	downloadPath := args[len(args)-1]

	if !initValidations(url) {
		return
	}

	clientId := soundcloud.GetClientId(url)

	if clientId == "" {
		fmt.Println("Something went wrong while getting the Client Id!")
		return
	}

	soundData := soundcloud.GetSoundMetaData(url, clientId)
	if soundData == nil {
		fmt.Printf("%s URL : %s doesn't return a valid track. Track is publicly accessed ?", theme.Red("[+]"), theme.Magenta(url))
		return
	}

	fmt.Printf("%s Track found. Title : %s - Duration : %s\n", theme.Green("[+]"), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))

	list := soundcloud.GetFormattedDL(soundData.Transcodes.Transcodings, clientId)

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
