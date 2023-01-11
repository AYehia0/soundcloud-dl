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

// download the track
func Download(track DownloadTrack, dlpath string) {
	// TODO: Prompt Y/N if the file exists and rename by adding _<random/date>.<ext>
	testPath := path.Join(dlpath, track.SoundData.Title+track.Ext)
	path, err := expandPath(testPath)

	// TODO: handle all different kind of errors
	if fileExists(path) || err != nil {
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
