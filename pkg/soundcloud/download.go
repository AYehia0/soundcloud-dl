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
	"strconv"

	"github.com/AYehia0/soundcloud-dl/pkg/theme"
	m "github.com/grafov/m3u8"
	"github.com/vbauerster/mpb/v8"
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
func downloadSeg(segmentURI string, file *os.File, bar *mpb.Bar) {
	resp, err := http.Get(segmentURI)
	var reader io.ReadCloser
	if err != nil {
		return
	}
	if bar != nil {
		reader = bar.ProxyReader(resp.Body)
	} else {
		reader = resp.Body
	}
	defer resp.Body.Close()
	defer reader.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return
	}
}

// parse the m3u8 to get into the playlist and extract the segments
// In our case, we're only interested in URLs of the segments
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

// making more requests to get the ContentLength of all the segments inside the playlist
// probablly it's not the best thing to do, but it's a nice to have the total size
// making a bar with 0 size making it really hard to :
// 1- Finish the bar to the end.
// 2- ProxyReader, when used displays wrong file sizes
// 3- Have to manually increment and SetTotal sizes, ahhhhh!
func addSegmentSizes(segments []string) int64 {
	var totalSize int64
	for _, segment := range segments {
		resp, _ := http.Head(segment)
		size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		totalSize += int64(size)
	}
	return totalSize
}

// using the goroutine to download each segment concurrently and wait till all finished
func DownloadM3u8(filepath string, segments []string, prog *mpb.Progress) error {

	file, _ := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	var bar *mpb.Bar = nil

	// getting the total size of all the segments = one track
	totalSize := addSegmentSizes(segments)

	if prog != nil {
		bar = theme.NewBar(prog, totalSize)
	}

	for _, segment := range segments {
		downloadSeg(segment, file, bar)
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

// a progress bar is passed which is used to display the progress bar (damn I am smart boi)
// if the download is done fine, the path of the created file is passed, else empty string is passed
// a better way is to return error, which i am going to do in the refactoring phase.
func Download(track DownloadTrack, dlpath string, prog *mpb.Progress) string {
	trackName := track.SoundData.Title + "[" + track.Quality + "]." + track.Ext
	path := validateDownload(dlpath, trackName)
	var reader io.ReadCloser

	// check if the track is hls
	if track.Quality != "low" {
		resp, err := http.Get(track.Url)
		if err != nil {
			return ""
		}
		defer resp.Body.Close()
		segments := getSegments(resp.Body)
		DownloadM3u8(path, segments, prog)

		return path
	}
	resp, err := http.Get(track.Url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	f, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if prog != nil {
		bar := theme.NewBar(prog, resp.ContentLength)
		reader = bar.ProxyReader(resp.Body)
	} else {
		reader = resp.Body
	}
	io.Copy(f, reader)

	return path
}
