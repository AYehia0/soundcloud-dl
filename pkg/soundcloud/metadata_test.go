package soundcloud_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
	"github.com/bogem/id3v2"
)

func TestAddMetaData(t *testing.T) {

	path := "../../testdata/"
	downloadTrack := soundcloud.DownloadTrack{
		Url:     "",
		Quality: "medium",
		Ext:     "mp3",
		SoundData: &soundcloud.SoundData{
			Id:           75825848,
			CreatedAt:    "2013-01-21T13:14:46Z",
			PermalinkUrl: "https://soundcloud.com/sobhi-mohamed5/99-118-mp4",
			ArtworkUrl:   "https://i1.sndcdn.com/artworks-000038806208-j0yb4i-large.jpg",
			Title:        "Test1",
			LikesCount:   355,
		},
	}

	// reading the img as [] byte
	imgBytes := readTestFile("artworks.jpg")
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		// return the url
		res.Write(imgBytes)
	}))

	// setting the url
	downloadTrack.SoundData.ArtworkUrl = testServer.URL

	filePath := filepath.Join(path, "Test[low]-nometa.mp3")
	soundcloud.AddMetadata(downloadTrack, filePath)

	// checking the tags
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})

	if err != nil {
		t.Errorf("Something went wrong reading the file tags : %s", err)
		return
	}

	defer tag.Close()

	if tag.Genre() != downloadTrack.SoundData.Genre {
		t.Errorf("Expected Genere to be %s, got %s", downloadTrack.SoundData.Genre, tag.Genre())
	}

	if tag.Title() != downloadTrack.SoundData.Title {
		t.Errorf("Expected Title to be %s, got %s", downloadTrack.SoundData.Title, tag.Title())
	}

	if tag.Year() != downloadTrack.SoundData.CreatedAt {
		t.Errorf("Expected Date to be %s, got %s", downloadTrack.SoundData.CreatedAt, tag.Year())
	}

	artist := strings.Split(downloadTrack.SoundData.PermalinkUrl, "/")[3]
	if tag.Artist() != artist {
		t.Errorf("Expected Artist to be %s, got %s", artist, tag.Artist())
	}

	// f := tag.GetFrames("APIC")[0]
	//
	// x, _ := f.(id3v2.PopularimeterFrame)
	pictures := tag.GetFrames(tag.CommonID("Attached picture"))
	for _, f := range pictures {
		pic, ok := f.(id3v2.PictureFrame)
		if !ok {
			log.Fatal("Couldn't assert picture frame")
		}

		mimeType := "image/jpeg"
		if pic.MimeType != mimeType {
			t.Errorf("Expected Mimetype to be %s, got %s", mimeType, pic.MimeType)
		}
		if !bytes.Equal(pic.Picture, imgBytes) {
			t.Errorf("Expected the picture to be the same as this : %s", "artworks.jpg")
		}
	}
	// // remove the tags
	tag.DeleteAllFrames()

	if err = tag.Save(); err != nil {
		t.Errorf("Something went wrong saving the file after removing the tags : %s", err)
		return
	}

}
