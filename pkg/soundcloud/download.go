// download the soundcloud tracks
package soundcloud

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"

	m "github.com/grafov/m3u8"
	bar "github.com/schollz/progressbar/v3"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// expand the given path ~/Desktop to the current logged in user /home/<username>/Desktop
func expandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil

}

// validate the given path
// and check if the file already exists or not
func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// extract the urls of the individual segment and then steam/download.
func downloadSeg(segmentURI string, file *os.File, dlbar *mpb.Bar) {
	resp, err := http.Get(segmentURI)

	if err != nil {
		return
	}

	reader := dlbar.ProxyReader(resp.Body)
	defer resp.Body.Close()

	defer reader.Close()

	dlbar.SetTotal(dlbar.Current()+resp.ContentLength, false)
	dlbar.IncrBy(int(resp.ContentLength))

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return
	}

}

func getSegments(body io.Reader) []string {
	segments := make([]string, 0)
	pl, listType, err := m.DecodeFrom(body, true)

	if err != nil {
		return nil
	}

	switch listType {
	case m.MEDIA:
		mediapl := pl.(*m.MediaPlaylist)
		for _, segment := range mediapl.Segments {
			if segment == nil {
				continue
			}
			segments = append(segments, segment.URI)
		}
	}
	return segments
}

// using the goroutine to download each segment concurrently and wait till all finished
func DownloadM3u8(filepath string, prog *mpb.Progress, segments []string) error {

	file, _ := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	// getting the total size of all the segments = one track
	// var totalSize int64
	// for _, segment := range segments {
	// 	resp, _ := http.Head(segment)
	// 	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	// 	totalSize += int64(size)
	// }

	dlbar := prog.AddBar(0,
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" -- "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	for _, segment := range segments {
		downloadSeg(segment, file, dlbar)
	}

	return nil
}

// before download validation
// return the path if everything is alright.
func validateDownload(dlpath string, trackName string) string {

	testPath := path.Join(dlpath, trackName)
	path, err := expandPath(testPath)

	// TODO: handle all different kind of errors
	if fileExists(path) || err != nil {
		return ""
	}
	return path
}

// download the track
func Download(track DownloadTrack, dlpath string, prog *mpb.Progress) string {
	// TODO: Prompt Y/N if the file exists and rename by adding _<random/date>.<ext>
	trackName := track.SoundData.Title + "[" + track.Quality + "]." + track.Ext
	path := validateDownload(dlpath, trackName)

	// check if the track is hls
	if track.Quality != "low" {

		resp, err := http.Get(track.Url)
		if err != nil {
			return ""
		}
		defer resp.Body.Close()

		segments := getSegments(resp.Body)
		DownloadM3u8(path, prog, segments)

		return path
	}
	resp, err := http.Get(track.Url)

	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// check if the file exists
	f, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := bar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)

	return path
}
