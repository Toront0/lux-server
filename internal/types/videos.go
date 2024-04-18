package types

import "time"

type VideoPreview struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Title string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	AuthorID string `json:"authorId"`
	Url string `json:"url"`
	ViewsAmount int `json:"viewsAmount"`
	AuthorFName string `json:"authorFName"`
	AuthorLName string `json:"authorLName"`
	AuthorPImg string `json:"authorPImg"`
}

type Video struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Title string `json:"title"`
	AuthorID string `json:"authorId"`
	Url string `json:"url"`
	Description string `json:"description"`
	ViewsAmount int `json:"viewsAmount"`
	AuthorFName string `json:"authorFName"`
	AuthorLName string `json:"authorLName"`
	AuthorPImg string `json:"authorPImg"`
}