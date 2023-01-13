// adding tags to the track after downloading it.
package soundcloud

import (
	"net/http"
	"strings"

	"github.com/AYehia0/soundcloud-dl/pkg/client"
	"github.com/bogem/id3v2"
)

func AddMetadata(track DownloadTrack, filePath string) error {
	t500 := "t500x500" // for getting a higher res img
	imgBytes := make([]byte, 0)

	// check for artist thing
	if track.SoundData.ArtworkUrl != "" {
		url := strings.Replace(track.SoundData.ArtworkUrl, "large", t500, 1)

		// fetching the data
		statusCode, data, err := client.Get(url)
		if err != nil || statusCode != http.StatusOK {
			return err
		}
		imgBytes = data
	}

	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return err
	}
	defer tag.Close()

	// setting metadata
	tag.SetTitle(track.SoundData.Title)
	tag.SetGenre(track.SoundData.Genre)
	tag.SetYear(track.SoundData.CreatedAt)

	// extracting the usr
	artistName := strings.Split(track.SoundData.PermalinkUrl, "/")
	tag.SetArtist(artistName[3])

	if imgBytes != nil {
		tag.AddAttachedPicture(
			id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    "image/jpeg",
				Picture:     imgBytes,
				Description: track.SoundData.Description, // well, coz why not :D
			},
		)
	}
	if err = tag.Save(); err != nil {
		return err
	}
	return nil

}
