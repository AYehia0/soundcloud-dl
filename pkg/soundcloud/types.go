package soundcloud

type Transcode struct {
	ApiUrl  string `json:"url"`
	Quality string `json:"quality"`
	Format  format `json:"format"`
}

type format struct {
	Protocol string `json:"protocol"`
	MimeType string `json:"mime_type"`
}

type SoundData struct {
	Id           int64      `json:"id"`
	Title        string     `json:"title"`
	CreatedAt    string     `json:"created_at"`
	Duration     int64      `json:"duration"`
	Kind         string     `json:"kind"`
	PermalinkUrl string     `json:"permalink_url"`
	UserId       int64      `json:"user_id"`
	ArtworkUrl   string     `json:"artwork_url"`
	Genre        string     `json:"genre"`
	Transcodes   Transcodes `json:"media"`
	LikesCount   int        `json:"likes_count"`
	Downloadable bool       `json:"downloadable"`
	Description  string     `json:"description,omitempty"`
}

type Transcodes struct {
	Transcodings []Transcode `json:"transcodings"`
}

type Media struct {
	Url string `json:"url"`
}

type DownloadTrack struct {
	Url       string
	Size      int
	Data      []byte
	Quality   string
	Ext       string
	SoundData *SoundData
}

type SearchResult struct {
	Sounds []SoundData `json:"collection"`
	Next   string      `json:"next_href"`
}
