// download the soundcloud tracks
package soundcloud

import (
	"io"
	"net/http"
	"os"

	bar "github.com/schollz/progressbar/v3"
)

// download the track
func Download(track DownloadTrack) {
	resp, err := http.Get(track.Url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// check if the file exists
	f, _ := os.OpenFile(track.SoundData.Title+track.Ext, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := bar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
