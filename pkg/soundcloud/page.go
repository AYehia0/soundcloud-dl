package soundcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/AYehia0/soundcloud-dl/pkg/client"
	"github.com/PuerkitoBio/goquery"
)

var Sound *SoundData

// convert the data map[string]any to well defined struct
// by marshalling then unmarshalling, which is probably isn't the best thing todo.
func formatData(data map[string]any) (*SoundData, error) {
	sound := &SoundData{}
	jsonBody, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonBody, &sound); err != nil {
		// do error check
		return nil, err
	}
	// adding the transcodings
	// go goes crazy mode lol
	transcodings := data["media"].(map[string]interface{})["transcodings"].([]any)

	for _, val := range transcodings {
		x := val.(map[string]any)
		encoding := Transcode{
			ApiUrl:  x["url"].(string),
			Quality: x["quality"].(string),
			Format:  x["format"].(map[string]interface{})["mime_type"].(string),
		}
		sound.Transcodings = append(sound.Transcodings, encoding)
	}

	return sound, nil
}

// extract some meta data under : window.__sc_hydration
func GetSoundMetaData(doc *goquery.Document) *SoundData {

	page := doc.First().Text()

	// TODO: You can do it in one shot, without splitting the text, but meh
	re := regexp.MustCompile(`(?m)^window.__sc_hydration = (.*?);\s*$`)
	matches := re.FindAllStringSubmatch(page, 1)

	if len(matches) < 1 {
		return nil
	}

	var result []map[string]interface{}
	data := strings.Split(matches[0][1], "\n")[0]

	json.Unmarshal([]byte(data), &result)

	// extracting the data
	// userData := result[6]["data"]
	soundData := result[7]["data"].(map[string]any)

	// TODO: return the error here
	Sound, _ = formatData(soundData)

	return Sound
}

func GetClientId(doc *goquery.Document) string {

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	// search for the client_id
	pattern := ",client_id:\"([^\"]*?.[^\"]*?)\""
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(body), 1)

	return matches[0][1]
}

func GetFormattedDL(data []Transcode, clientId string) []DownloadTrack {
	ext := "mp3" // the default extension type
	tracks := make([]DownloadTrack, 0)

	for _, tcode := range data {
		url := tcode.ApiUrl + "?client_id=" + clientId
		statusCode, body, err := client.Get(url)
		if err != nil && statusCode != http.StatusOK {
			continue
		}
		q := mapQuality(tcode.ApiUrl, tcode.Format)
		if q == "ogg" {
			ext = q
		}
		mediaUrl := Media{}
		json.Unmarshal(body, &mediaUrl)

		tmpTrack := DownloadTrack{
			Url:       mediaUrl.Url,
			Quality:   q,
			SoundData: Sound,
			Ext:       ext,
		}
		tracks = append(tracks, tmpTrack)
	}
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
