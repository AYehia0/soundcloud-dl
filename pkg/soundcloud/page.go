package soundcloud

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
	dataSound, _ := formatData(soundData)

	return dataSound
}
