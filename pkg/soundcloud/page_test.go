package soundcloud_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
)

var (
	trackUrl = "https://soundcloud.com/sobhi-mohamed5/99-118-mp4"
	clientId = "ZQvaVYuPpe0Pg7Ga7V24qFseYl6eTK73"
)

// read the html file
func readTestFile(testfile string) []byte {
	testPath := "../../testdata/"
	content, err := ioutil.ReadFile(filepath.Join(testPath, testfile))

	if err != nil {
		return nil
	}

	return content
}

// TODO: Test if server fails to respond
func TestGetClientId(t *testing.T) {
	expectedData := "ZQvaVYuPpe0Pg7Ga7V24qFseYl6eTK73"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(readTestFile("user-soundcloud.html"))
	}))

	defer testServer.Close()

	clientId := soundcloud.GetClientId(testServer.URL)

	if clientId != expectedData {
		t.Errorf("Something wen't wrong, expected clientId to be %s", expectedData)
	}
}

func TestGetSoundMetaData(t *testing.T) {

	// the expected data isn't complete as I am lazy and not goint to check for everything
	expectedData := &soundcloud.SoundData{
		Id:           75825848,
		CreatedAt:    "2013-01-21T13:14:46Z",
		PermalinkUrl: "https://soundcloud.com/sobhi-mohamed5/99-118-mp4",
		ArtworkUrl:   "https://i1.sndcdn.com/artworks-000038806208-j0yb4i-large.jpg",
		LikesCount:   355,
	}
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(readTestFile("track-info.json"))
	}))

	defer testServer.Close()

	data := soundcloud.GetSoundMetaData(testServer.URL, trackUrl, clientId)

	if data == nil {
		t.Errorf("Expected data of type %T, got : %v", expectedData, data)
		return
	}

	if data.Id != expectedData.Id {
		t.Errorf("Expected Id, to be %v, got %v", expectedData.Id, data.Id)
	}

	if data.LikesCount != expectedData.LikesCount {
		t.Errorf("Expected LikesCount, to be %v, got %v", expectedData.LikesCount, data.LikesCount)
	}

	if data.PermalinkUrl != expectedData.PermalinkUrl {
		t.Errorf("Expected PermalinkUrl, to be %s, got %s", expectedData.PermalinkUrl, data.PermalinkUrl)
	}

	if data.ArtworkUrl != expectedData.ArtworkUrl {
		t.Errorf("Expected ArtworkUrl, to be %s, got %s", expectedData.ArtworkUrl, data.ArtworkUrl)
	}
}

func TestGetFormattedDL(t *testing.T) {

	mediaUrl := "https://cf-media.sndcdn.com/nDzLF7r2P2Xd.128.mp3?Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiKjovL2NmLW1lZGlhLnNuZGNkbi5jb20vbkR6TEY3cjJQMlhkLjEyOC5tcDMqIiwiQ29uZGl0aW9uIjp7IkRhdGVMZXNzVGhhbiI6eyJBV1M6RXBvY2hUaW1lIjoxNjczNzY3Mzk1fX19XX0_&Signature=LkOk4zO9a86Hnr62b0CNG~6RqZg-o70Z3Xun93pxGMt3ntZfWU0WNjMqIvjLEALXyr6f3zsGoGXdUIU4j92ObbkCZjxWNXt3p-EMkKpSATqrc48WD7OwjFr-7-xW4eFzF6SgCYG3AUQJghzv8vBGcMzV7ZCZPmsb8M88vP~UiO3K-7jyT9UJt6jWzuGGS6WDzyNoVkg7yhhuPg-m8i3gJHNZR-lN7190xP11bCF8qZUZVq-ymXmNxovzhBnO9HbzRTviGPNFfpLLKs5HPMGEiK3HKsoSTVukUTfvn-PInyZWAGZDZmvPChIYkRSD1S51qLvsyBZkExvVcpvxebZ-5A__&Key-Pair-Id=APKAI6TU7MMXM5DG6EPQ"
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		// return the url
		res.Write(readTestFile("media-url.json"))
	}))

	defer testServer.Close()

	/*
		[{
		https://api-widget.soundcloud.com/media/soundcloud:tracks:75825848/7335019e-2579-47b2-bdc4-5ec031ed63ae/stream/hls
		sq
		{hls audio/mpeg}}

		{https://api-widget.soundcloud.com/media/soundcloud:tracks:75825848/7335019e-2579-47b2-bdc4-5ec031ed63ae/stream/progressive
		sq
		{progressive audio/mpeg}
		}]
	*/
	// here I am testing on just a single Transcodings instance
	// pretty ugly testing but it works :D
	track := &soundcloud.SoundData{
		Id:           75825848,
		CreatedAt:    "2013-01-21T13:14:46Z",
		PermalinkUrl: "https://soundcloud.com/sobhi-mohamed5/99-118-mp4",
		ArtworkUrl:   "https://i1.sndcdn.com/artworks-000038806208-j0yb4i-large.jpg",
		Transcodes: soundcloud.Transcodes{
			Transcodings: []soundcloud.Transcode{
				{ApiUrl: testServer.URL, Quality: "sq", Format: soundcloud.Format{
					Protocol: "hls",
					MimeType: "audio/mpeg",
				}},
			},
		},
	}

	formattedDownloadT := soundcloud.GetFormattedDL(track, clientId)

	// the expected length is 1
	if len(formattedDownloadT) < 0 {
		t.Errorf("Expected the length of %T to be >= 1, got : %v", formattedDownloadT, len(formattedDownloadT))
		return
	}

	if formattedDownloadT[0].SoundData.Id != track.Id {
		t.Errorf("Expected Id, to be %v, got %v", track.Id, formattedDownloadT[0].SoundData.Id)
	}

	if formattedDownloadT[0].Url != mediaUrl {
		t.Errorf("Expected MediaUrl to be %s, got : %s", mediaUrl, formattedDownloadT[0].Url)
	}
}

func TestSearchTracksByKeyword(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(readTestFile("search-result.json"))
	}))

	defer testServer.Close()

	keyword := "surah yasin"
	offset := 0
	searchResult := soundcloud.SearchTracksByKeyword(testServer.URL, keyword, offset, clientId)

	expectedIds := []int64{
		781025620,
		129803149,
		893465524,
		103207875,
		163687464,
		144965530,
	}

	ids := make([]int64, 0)
	for _, res := range searchResult.Sounds {
		ids = append(ids, res.Id)
	}

	if len(searchResult.Sounds) != len(expectedIds) {
		t.Errorf("Expected length of search result to be : %v, got %v", len(expectedIds), len(searchResult.Sounds))
	}

	if !reflect.DeepEqual(ids, expectedIds) {
		t.Errorf("Expected search Ids : %v to be equal, got %v", expectedIds, ids)
	}

}
