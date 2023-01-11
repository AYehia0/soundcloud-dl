// download the soundcloud tracks
package soundcloud

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sync"

	m "github.com/grafov/m3u8"
	bar "github.com/schollz/progressbar/v3"
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
	fmt.Printf("File: %s Exists\n", path)
	return true
}

// extract the urls of the individual segment and then steam/download.
func downloadSeg(wg *sync.WaitGroup, segment *m.MediaSegment, file *os.File, dlbar *bar.ProgressBar) {
	defer wg.Done()
	resp, err := http.Get(segment.URI)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	// append to the file
	_, err = io.Copy(io.MultiWriter(file, dlbar), resp.Body)

	if err != nil {
		return
	}

}

// using the goroutine to download each segment concurrently and wait till all finished
func downloadM3u8(m3u8Url string, filepath string) error {
	resp, err := http.Get(m3u8Url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	pl, listType, err := m.DecodeFrom(resp.Body, true)
	if err != nil {
		return err
	}

	dlbar := bar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	// create a file to add content to
	file, _ := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// file, _ := os.Create(filepath)

	// the go routine now
	var wg sync.WaitGroup
	switch listType {
	case m.MEDIA:
		mediapl := pl.(*m.MediaPlaylist)
		for _, segment := range mediapl.Segments {
			if segment == nil {
				continue
			}
			wg.Add(1)
			// the added go routine has nothing to do here, as `go` keyword causes some issues idk why
			downloadSeg(&wg, segment, file, dlbar)
		}
	default:
		return errors.New("Unsupported type!")
	}

	return nil
}

// download the track
func Download(track DownloadTrack, dlpath string) {
	// TODO: Prompt Y/N if the file exists and rename by adding _<random/date>.<ext>
	testPath := path.Join(dlpath, track.SoundData.Title+track.Ext)
	path, err := expandPath(testPath)

	// TODO: handle all different kind of errors
	if fileExists(path) || err != nil {
		return
	}

	// check if the track is hls
	if track.Quality != "low" {
		downloadM3u8(track.Url, path)
		return
	}
	resp, err := http.Get(track.Url)

	if err != nil {
		return
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
}
