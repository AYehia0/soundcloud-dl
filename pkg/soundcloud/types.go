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
	Transcodings []Transcode
	TrackAuth    string `json:"track_authorization"`
	LikesCount   int    `json:"likes_count"`
	Downloadable bool   `json:"downloadable"`
	Description  string `json:"description,omitempty"`
}
