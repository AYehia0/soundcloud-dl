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

func Sc(args []string, bestQuality bool) {

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
		fmt.Println("URL doesn't exist : status is : ", statusCode)
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
	qualities := getQualities(list)
	if !bestQuality {
		defaultQuality = chooseQuality(qualities)
	} else {
		defaultQuality = getHighestQuality(qualities)
	}

	track := chooseTrackDownload(list, defaultQuality)
	filePath := soundcloud.Download(track, downloadPath)

	// add tags
	err = soundcloud.AddMetadata(track, filePath)
	if err != nil {
		fmt.Printf("Error happend while adding tags to the track : %s\n", err)
	}
	fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(filePath))
}
