package types

import "time"

type CommunityPreview struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Category string `json:"category"`
	ProfileImg string `json:"profileImg"`
}

type Community struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Title int `json:"title"`
	Category string `json:"category"`
	ProfileImg string `json:"profileImg"`
	Description string `json:"description"`
}