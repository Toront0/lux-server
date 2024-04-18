package handlers

import (
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"

	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type musicHandler struct {
	store services.MusicStorer
	cld *cloudinary.Cloudinary
}

func NewMusicHandler(store services.MusicStorer, cld *cloudinary.Cloudinary) *musicHandler {

	return &musicHandler{
		store: store,
		cld: cld,
	} 
}

func (h *musicHandler) HandleGetSongs(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")

	userID := middleware.UserFromContext(r.Context())

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page query to get music %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetSongs(userID, _page)

	if err != nil {
		fmt.Printf("could not get songs %s", err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *musicHandler) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {

	res, err := h.store.GetPlaylists()

	if err != nil {
		fmt.Printf("could not get playlists %s", err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *musicHandler) HandleGetPlaylistDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	pID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not convert provided playlist ID %s", err)
		w.WriteHeader(400)
		return
	}

	fmt.Printf("playlist id id %d", pID)

	res, err := h.store.GetPlaylistDetail(pID)

	if err != nil {
		fmt.Printf("could not find playlist with such an id %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *musicHandler) HandleGetPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	pID := r.PathValue("playlistId")
	uID := r.PathValue("userId")

	page := r.URL.Query().Get("page")


	_pID, err := strconv.Atoi(pID)

	if err != nil {
		fmt.Printf("could not convert provided playlist ID %s", err)
		w.WriteHeader(400)
		return
	}

	_uID, err := strconv.Atoi(uID)

	if err != nil {
		fmt.Printf("could not convert user id to get playlist %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not provided page query for fetch playlist %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetPlaylistSongs(_pID, _uID, _page)

	if err != nil {
		fmt.Printf("could not find the playlist %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *musicHandler) HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Title string `json:"title"`
		CoverImg string `json:"coverImg"`
		Songs []int `json:"songs"`
	} {
		Title: "",
		CoverImg: "",
		Songs: []int{},
	}

	
	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	query := "insert into user_music_playlists "

	if req.CoverImg != "" {
		imgID := "playlist-" + req.Title + "-creator-id" + strconv.Itoa(userID)

		res, err := h.cld.Upload.Upload(context.Background(), req.CoverImg, uploader.UploadParams{PublicID: imgID, Folder: "social-media/playlist_cover"})

		if err != nil {
			fmt.Printf("could not upload playlist's cover img %s", err)
			w.WriteHeader(500)
			return
		}



		query += "(cover_img, title, user_id) values (" + "'" + res.SecureURL + "'" + ", " + "'" + req.Title + "'" + ", " + strconv.Itoa(userID) + ")" 

	} else {

		query += "(title, user_id) values (" + "'" + req.Title + "'" + ", " + strconv.Itoa(userID) + ")"

	}

	query += " returning id"

	id, err := h.store.CreatePlaylist(query)

	if err != nil {
		fmt.Printf("could not create playlist %s", err)
		w.WriteHeader(400)
		return
	}

	query = "update user_music set playlist_id = " + strconv.Itoa(id)

	for i, s := range req.Songs {

		if i == 0 {

			query += " where song_id = " + strconv.Itoa(s) + " and user_id =" + strconv.Itoa(userID)

		} else if i == len(req.Songs) - 1 {

			query += " or song_id = " + strconv.Itoa(s) + " and user_id =" + strconv.Itoa(userID)

		} else {

			query += " or song_id = " + strconv.Itoa(s) + " and user_id =" + strconv.Itoa(userID)

		}

	}
 

	err = h.store.UpdateSongsPlaylistID(query)
	
	if err != nil {
		fmt.Printf("could not change songs playlist id %s", err)
		w.WriteHeader(400)
		return
	}


}

func (h *musicHandler) HandleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_id, err := strconv.Atoi(id)

	userID := middleware.UserFromContext(r.Context())

	if err != nil {
		fmt.Printf("could not parse id to delete playlist %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.DeletePlaylist(_id, userID)

	if err != nil {
		fmt.Printf("could not delete playlist %s", err)
		w.WriteHeader(400)
		return
	}

}

func (h *musicHandler) HandleGetAvailableAndPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	pID := r.PathValue("playlistId")
	uID := r.PathValue("userId")

	page := r.URL.Query().Get("page")


	_pID, err := strconv.Atoi(pID)

	if err != nil {
		fmt.Printf("could not convert provided playlist ID %s", err)
		w.WriteHeader(400)
		return
	}

	_uID, err := strconv.Atoi(uID)

	if err != nil {
		fmt.Printf("could not convert user id to get playlist %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not provided page query for fetch playlist %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetAvailableAndPlaylistSongs(_pID, _uID, _page)

	if err != nil {
		fmt.Printf("could not fetch GetAvailableAndPlaylistSongs %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)


}

func (h *musicHandler) HandleAddSongToUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userID := middleware.UserFromContext(r.Context())

	_id, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not convert music id to add it to user %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.AddSongToUser(_id, userID)

	if err != nil {
		fmt.Printf("could not add the song to user %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *musicHandler) HandleDeleteUserSong(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userID := middleware.UserFromContext(r.Context())

	_id, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not convert music id to delete it %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.DeleteUserSong(_id, userID)

	if err != nil {
		fmt.Printf("could not delete user song %s", err)
		w.WriteHeader(400)
		return
	}
}