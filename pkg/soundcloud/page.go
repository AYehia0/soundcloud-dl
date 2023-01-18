package soundcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/AYehia0/soundcloud-dl/pkg/client"
	"github.com/PuerkitoBio/goquery"
)

var Sound *SoundData

// extract some meta data under : window.__sc_hydration
// check if the track exists and open to public
func GetSoundMetaData(apiUrl string, url string, clientId string) *SoundData {

	statusCode, body, err := client.Get(apiUrl)

	if err != nil || statusCode != http.StatusOK {
		return nil
	}

	json.Unmarshal(body, &Sound)

	return Sound
}

func GetClientId(url string) string {

	if url == "" {
		// the best url ever, if you find this then you're so cool :D I love you :DDD
		url = "https://soundcloud.com/ahmed-yehia0"
	}

	statusCode, bodyData, err := client.Get(url)

	if err != nil {
		log.Fatalf("An Error : %s happended while requesting : %s", err, url)
	}
	if statusCode != http.StatusOK {
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyData))

	// find the last src under the body
	apiurl, exists := doc.Find("body > script").Last().Attr("src")
	if !exists {
		return ""
	}

	// making a GET request to find the client_id
	resp, err := http.Get(apiurl)
	if err != nil {
		fmt.Printf("Something went wrong while requesting : %s , Error : %s", apiurl, err)
	}

	// reading the body
	bodyData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// search for the client_id
	pattern := ",client_id:\"([^\"]*?.[^\"]*?)\""
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(bodyData), 1)

	return matches[0][1]
}

func GetFormattedDL(track *SoundData, clientId string) []DownloadTrack {

	ext := "mp3" // the default extension type
	tracks := make([]DownloadTrack, 0)
	data := track.Transcodes.Transcodings
	var wg sync.WaitGroup

	for _, tcode := range data {
		wg.Add(1)
		go func(tcode Transcode) {
			defer wg.Done()

			url := tcode.ApiUrl + "?client_id=" + clientId
			statusCode, body, err := client.Get(url)
			if err != nil && statusCode != http.StatusOK {
				return
			}
			q := mapQuality(tcode.ApiUrl, tcode.Format.MimeType)
			if q == "high" {
				ext = "ogg"
			}
			mediaUrl := Media{}
			dec := json.NewDecoder(bytes.NewReader(body))
			if err := dec.Decode(&mediaUrl); err != nil {
				log.Println("Error decoding json: ", err)
				return
			}
			tmpTrack := DownloadTrack{
				Url:       mediaUrl.Url,
				Quality:   q,
				SoundData: track,
				Ext:       ext,
			}
			tracks = append(tracks, tmpTrack)

		}(tcode)
	}
	wg.Wait()
	return tracks
}

// check if the trackUrl is mp3:progressive or ogg:hls
func mapQuality(url string, format string) string {
	tmp := strings.Split(url, "/")
	if tmp[len(tmp)-1] == "hls" && strings.HasPrefix(format, "audio/ogg") {
		return "high"
	} else if tmp[len(tmp)-1] == "hls" && strings.HasPrefix(format, "audio/mpeg") {
		return "medium"
	}
	return "low"
}

func SearchTracksByKeyword(apiUrl string, keyword string, offset int, clientId string) *SearchResult {

	statusCode, body, err := client.Get(apiUrl)

	if err != nil && statusCode != http.StatusOK {
		return nil
	}

	var result = SearchResult{}
	json.Unmarshal(body, &result)

	return &result
}
