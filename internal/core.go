package internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/AYehia0/soundcloud-dl/pkg/theme"
	"github.com/vbauerster/mpb/v8"
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

		fmt.Printf("%s %s found. Title : %s - Duration : %s\n", theme.Green("[+]"), strings.Title(soundData.Kind), theme.Magenta(soundData.Title), theme.Magenta(theme.FormatTime(soundData.Duration)))
	}

	// check if the url is a playlist
	if soundData.Kind == "playlist" {
		var wg sync.WaitGroup
		plDownloadTracks := getPlaylistDownloadTracks(soundData, clientId)
		p := mpb.New(mpb.WithWaitGroup(&wg),
			mpb.WithWidth(64),
			mpb.WithRefreshRate(180*time.Millisecond),
		)

		for _, dlT := range plDownloadTracks {

			wg.Add(1)

			go func(dlT []soundcloud.DownloadTrack) {
				defer wg.Done()
				// bestQuality is true to avoid prompting the user for quality choosing each time and speed up
				// TODO: get a single progress bar, this will require the use of "https://github.com/cheggaaa/pb" since the current pb doesn't support download pool (I think)
				t := getTrack(dlT, true)
				fp := soundcloud.Download(t, downloadPath, p)

				// silent indication of already existing files
				if fp == "" {
					return
				}
				soundcloud.AddMetadata(t, fp)
			}(dlT)

		}
		wg.Wait()

		fmt.Printf("\n%s Playlist saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
		return
	}

	downloadTracks := soundcloud.GetFormattedDL(soundData, clientId)

	track := getTrack(downloadTracks, bestQuality)
	filePath := soundcloud.Download(track, downloadPath, nil)

	// add tags
	if filePath == "" {
		fmt.Printf("\n%s Track was already saved to : %s\n", theme.Green("[-]"), theme.Magenta(downloadPath))
		return
	}
	err := soundcloud.AddMetadata(track, filePath)
	if err != nil {
		fmt.Printf("Error happend while adding tags to the track : %s\n", err)
	}
	fmt.Printf("\n%s Track saved to : %s\n", theme.Green("[-]"), theme.Magenta(filePath))
}
