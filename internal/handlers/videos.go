package handlers

import (
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"

	"fmt"
	"net/http"
	"encoding/json"
	"strconv"

)

type videoHandler struct {
	store services.VideoStorer
}

func NewVideoHandler(store services.VideoStorer) *videoHandler {
	
	return &videoHandler{
		store: store,
	}
}

func (h *videoHandler) HandleGetVideos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("invalid page was provided to get comments %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetVideos(_page)

	if err != nil {
		fmt.Printf("could not get videos %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *videoHandler) HandleGetVideo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	vID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not convert video id %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetVideo(vID)

	if err != nil {
		fmt.Printf("could not get video %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *videoHandler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	page := r.URL.Query().Get("page")

	userID := middleware.UserFromContext(r.Context())

	pID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("invalid video id was provided to get comments %s", err)
		w.WriteHeader(400)
		return
	}

	p, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("invalid page was provided to get comments %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetComments(pID, p, userID)

	if err != nil {
		fmt.Printf("could not get video comments %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *videoHandler) HandleLikeVideo(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		VideoID int `json:"videoId"`
	} {
		VideoID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.LikeVideo(req.VideoID, userID)

	if err != nil {
		fmt.Printf("could not like the post %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *videoHandler) HandleDeleteLike(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		VideoID int `json:"videoId"`
	} {
		VideoID: 0,
	}

	userID := middleware.UserFromContext(r.Context())

	json.NewDecoder(r.Body).Decode(req)

	err := h.store.DeleteLike(req.VideoID, userID)

	if err != nil {
		fmt.Printf("could not delete like from the post %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *videoHandler) HandleInsertComment(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		VideoID int `json:"videoId"`
		Content string `json:"content"`
	} {
		VideoID: 0,
		Content: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.InsertComment(req.VideoID, userID, req.Content)

	if err != nil {
		fmt.Printf("could not insert video comment %s", err)
		w.WriteHeader(400)
		return
	}

}


func (h *videoHandler) HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.LikeComment(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not like the post comment %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *videoHandler) HandleDeleteLikeComment(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.DeleteLikeComment(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not delete like from the post comment %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *videoHandler) HandleGetCommentReplies(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	userID := middleware.UserFromContext(r.Context())

	cID, err := strconv.Atoi(id)


	if err != nil {
		fmt.Printf("invalid video comment id was provided to get comment replies %s", err)
		w.WriteHeader(400)
		return
	}

	l, err := strconv.Atoi(limit)

	if err != nil {
		fmt.Printf("invalid limit url query was provided to get comment replies %s", err)
		w.WriteHeader(400)
		return
	}

	o, err := strconv.Atoi(offset)

	if err != nil {
		fmt.Printf("invalid offset url query was provided to get comment replies %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetCommentReplies(cID, l, o, userID)

	if err != nil {
		fmt.Printf("could not get video comment replies %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *videoHandler) HandleInsertCommentReply(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		VideoID int `json:"videoId"`
		Content string `json:"content"`
		CommentID int `json:"commentId"`
	} {
		VideoID: 0,
		Content: "",
		CommentID: 0,
	}

	userID := middleware.UserFromContext(r.Context())

	json.NewDecoder(r.Body).Decode(req)

	err := h.store.InsertCommentReply(userID, req.CommentID, req.VideoID, req.Content)

	if err != nil {
		fmt.Printf("could not insert video comment reply %s", err)
		w.WriteHeader(400)
		return
	}

}

func (h *videoHandler) HandleLikeCommentReply(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.LikeCommentReply(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not like the video comment reply %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *videoHandler) HandleDeleteCommentLikeReply(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.DeleteCommentLikeReply(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not delete like from the video comment reply %s", err)
		w.WriteHeader(400)
		return
	}
}