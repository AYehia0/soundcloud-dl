// get list of DownloadTrack of all the urls in the playlist
// with the number of urls in the playlist

package soundcloud

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/AYehia0/soundcloud-dl/pkg/client"
)

type Id struct {
	Id int `json:"id"`
}

type Tracktmp struct {
	Tracks []Id `json:"tracks"`
}

// get all the urls inside the playlist
func GetPlaylistTracks(track *SoundData, clientId string) []SoundData {
	ids := make([]string, 0)
	trackIds := Tracktmp{}
	plApiUrl := GetTrackInfoAPIUrl(track.PermalinkUrl, clientId)

	statusCode, data, err := client.Get(plApiUrl)

	if err != nil || statusCode != http.StatusOK {
		return nil
	}

	//fmt.Println(string(data))
	json.Unmarshal(data, &trackIds)

	for _, t := range trackIds.Tracks {
		ids = append(ids, strconv.Itoa(t.Id))
	}

	tApiUrl := GetTracksByIdsApiUrl(ids, clientId)

	statusCode, data, err = client.Get(tApiUrl)

	if err != nil || statusCode != http.StatusOK {
		return nil
	}

	sounds := make([]SoundData, 0)

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&sounds); err != nil {
		log.Println("Error decoding json: ", err)
		return nil
	}

	return sounds
}
