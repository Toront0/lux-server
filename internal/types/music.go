package types

type MusicPlaylist struct {
	ID int `json:"id"`
	Title string `json:"title"`
	CoverImg string `json:"coverImg"`
	CreatorID int `json:"creatorId"`
	CreatorFName string `json:"creatorFName"`
	CreatorLName string `json:"creatorLName"`
	CreatorPImg string `json:"creatorPImg"`
} 

type MusicPlaylistPreview struct {
	ID int `json:"id"`
	Title string `json:"title"`
	CoverImg *string `json:"coverImg"`
	CreatorID int `json:"creatorId"`
} 

type Song struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Performer string `json:"performer"`
	Cover string `json:"cover"`
	Url string `json:"url"`
	// 0 if does not exist otherwise > 0
	IsInMyList *int `json:"isInMyList"`
}

type PlaylistSong struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Performer string `json:"performer"`
	Cover string `json:"cover"`
	Url string `json:"url"`
	PlaylistID *int `json:"playlistId"`
}