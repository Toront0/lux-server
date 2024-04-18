package types

import (
	"time"
)

type CommentReply struct {
	ID int `json:"id"`
	CreatedAt any `json:"createdAt"`
	Content string `json:"content"`
	AuthorID int `json:"authorID"`
	AuthorFName string `json:"authorFName"`
	AuthorLName string `json:"authorLName"`
	AuthorPImg string `json:"authorPImg"`
	LikesAmount *int `json:"likesAmount"`
	IsRequesterLiked *int `json:"isRequesterLiked"`
}

type Comment struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Content string `json:"content"`
	AuthorID int `json:"authorID"`
	AuthorFName string `json:"authorFName"`
	AuthorLName string `json:"authorLName"`
	AuthorPImg string `json:"authorPImg"`
	LikesAmount int `json:"likesAmount"`
	RepliesAmount int `json:"repliesAmount"`
	IsRequesterLiked int `json:"isRequesterLiked"`
}

type Post struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Content string `json:"content"`
	AuthorID int `json:"authorId"`
	AuthorName string `json:"authorName"`
	AuthorPImg string `json:"authorPImg"`
	LikesAmount int `json:"likesAmount"`
	IsRequesterLiked int `json:"isRequesterLiked"`
	CommentsAmount int `json:"commentsAmount"`
	PostMedia []string `json:"postMedia"`
}

type PostResponse struct {
	Count int `json:"count"`
	Posts []Post `json:"posts"`
}

type PostResponse2 struct {
	Posts []Post `json:"posts"`
	
}