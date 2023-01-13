package soundcloud

type Transcode struct {
	ApiUrl  string
	Quality string
	Format  string
}

type SoundData struct {
	Id           int64  `json:"id"`
	Title        string `json:"title"`
	CreatedAt    string `json:"created_at"`
	Duration     int64  `json:"duration"`
	Kind         string `json:"kind"`
	PermalinkUrl string `json:"permalink_url"`
	UserId       int64  `json:"user_id"`
	ArtworkUrl   string `json:"artwork_url"`
	Genre        string `json:"genre"`
	Transcodings []Transcode
	TrackAuth    string `json:"track_authorization"`
	LikesCount   int    `json:"likes_count"`
	Downloadable bool   `json:"downloadable"`
	Description  string `json:"description,omitempty"`
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
