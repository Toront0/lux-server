package handlers

import (
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"

	"github.com/cloudinary/cloudinary-go/v2"
	// "github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"


	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"context"

)

type postHandler struct {
	store services.PostStorer
	cld *cloudinary.Cloudinary
}

func NewPostHandler(store services.PostStorer, cld *cloudinary.Cloudinary) *postHandler {

	return &postHandler{
		store: store,
		cld: cld,
	}

}

func (h *postHandler) HandleGetPosts(w http.ResponseWriter, r *http.Request) {

	userID := middleware.UserFromContext(r.Context())
	
	req := r.URL.Query().Get("page")

	page, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("could not parse page to get posts %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetPosts(page, userID)

	if err != nil {
		fmt.Printf("could not get posts %s", err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *postHandler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userID := middleware.UserFromContext(r.Context())

	pID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse id of requested post %s", err)
		w.WriteHeader(400)	
		return
	}

	res, err := h.store.GetPost(pID, userID)

	if err != nil {
		fmt.Printf("could not get post %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *postHandler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	page := r.URL.Query().Get("page")

	userID := middleware.UserFromContext(r.Context())

	pID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("invalid post id was provided to get comments %s", err)
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
		fmt.Printf("could not get post comments %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *postHandler) HandleInsertComment(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		PostID int `json:"postId"`
		Content string `json:"content"`
	} {
		PostID: 0,
		Content: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.InsertComment(req.PostID, userID, req.Content)

	if err != nil {
		fmt.Printf("could not insert post comment %s", err)
		w.WriteHeader(400)
		return
	}



}

func (h *postHandler) HandleGetCommentReplies(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	userID := middleware.UserFromContext(r.Context())


	cID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("invalid post comment id was provided to get comment replies %s", err)
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
		fmt.Printf("could not get post comment replies %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *postHandler) HandleInsertCommentReply(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		PostID int `json:"postId"`
		Content string `json:"content"`
		CommentID int `json:"commentId"`
	} {
		PostID: 0,
		Content: "",
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.InsertCommentReply(userID, req.CommentID, req.PostID, req.Content)

	if err != nil {
		fmt.Printf("could not insert comment reply %s", err)
		w.WriteHeader(400)
		return
	}

}

func (h *postHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		PostID int `json:"postId"`
	} {
		PostID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.LikePost(req.PostID, userID)

	if err != nil {
		fmt.Printf("could not like the post %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *postHandler) HandleDeleteLike(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		PostID int `json:"postId"`
	} {
		PostID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.DeleteLike(req.PostID, userID)

	if err != nil {
		fmt.Printf("could not delete like from the post %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *postHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Content string `json:"content"`
		Imgs []string `json:"imgs"`
	} {
		Content: "",
		Imgs: []string{},
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	imgsUrls := []string{}


	if len(req.Imgs) > 0 {
		
		for i, img := range req.Imgs {

			imgID := strconv.Itoa(i + 1) + "-image-of-" + strconv.Itoa(userID) + "-post"

			res, err := h.cld.Upload.Upload(context.Background(), img, uploader.UploadParams{PublicID: imgID, Folder: "social-media/posts"});

			

			if err != nil {
				fmt.Printf("could not upload img to the cloudinary %s", err)
				return
			}

			fmt.Println("resp.SecureURL:", res.SecureURL)

			imgsUrls = append(imgsUrls, res.SecureURL)

		}

	}

	err := h.store.CreatePost(req.Content, userID, imgsUrls)

	if err != nil {
		fmt.Printf("could not create the post %s", err)
		w.WriteHeader(400)
		return
	}

}

func (h *postHandler) HandleLikeComment(w http.ResponseWriter, r *http.Request) {
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

func (h *postHandler) HandleDeleteLikeComment(w http.ResponseWriter, r *http.Request) {
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

func (h *postHandler) HandleLikeCommentReply(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.LikeCommentReply(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not like the post comment reply %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *postHandler) HandleDeleteCommentLikeReply(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		CommentID int `json:"commentId"`
	} {
		CommentID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.DeleteCommentLikeReply(req.CommentID, userID)

	if err != nil {
		fmt.Printf("could not delete like from the post comment reply %s", err)
		w.WriteHeader(400)
		return
	}
}