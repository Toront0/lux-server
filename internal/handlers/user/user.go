package user

import (
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"
	"guthub.com/Toront0/lux-server/internal/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"fmt"
	"context"
	"net/http"
	"encoding/json"
	"strconv"
	"log"
	"time"

)

type userHandler struct {
	store services.UserStorer
	hub *Hub
	cld *cloudinary.Cloudinary
}

func NewUserHandler(store services.UserStorer, cld *cloudinary.Cloudinary) *userHandler {

	hub := NewHub()

	go hub.Run()

	return &userHandler{
		store: store,
		hub: hub,
		cld: cld,
	}
}

func (h *userHandler) HandleGetUserDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userID := middleware.UserFromContext(r.Context())

	uID, err := strconv.Atoi(id)


	if err != nil {
		fmt.Printf("could not convert user ID %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserDetail(uID, userID)

	if err != nil {
		fmt.Printf("could not get user Detail %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userHandler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	page := r.URL.Query().Get("page")

	uID, err := strconv.Atoi(id)


	if err != nil {
		fmt.Printf("could not get user posts %s", err)
		w.WriteHeader(400)
		return
	}

	

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get friends %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserPosts(uID, _page)

	if err != nil {
		fmt.Printf("could not get user posts %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *userHandler) HandleGetUserVideos(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") 
	page := r.URL.Query().Get("page")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user id to get videos %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get friends %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetUserVideos(uID, _page)

	if err != nil {
		fmt.Printf("could not get user videos %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userHandler) HandleGetUserFriends(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") 
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user Id to get friends %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get friends %s", err)
		w.WriteHeader(400)
		return
	}


	res, err := h.store.GetUserFriends(uID, _page, search)

	if err != nil {
		fmt.Printf("could not get user friends %s", err)
		w.WriteHeader(404)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userHandler) HandleGetUserMusic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	page := r.URL.Query().Get("page")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse id to get user communities %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get user communities %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserMusic(uID, _page)

	if err != nil {
		fmt.Printf("could not get user music %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userHandler) HandleGetUserMusicPlaylists(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse id to get user communities %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserPlaylists(uID)

	if err != nil {
		fmt.Printf("could not get user music %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)

}

func (h *userHandler) HandleGetUserFollowers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") 
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user Id to get followers %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get followers %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserFollowers(uID, _page, search)

	if err != nil {
		fmt.Printf("could not get user followers %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *userHandler) HandleGetUserFollowings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") 
	page := r.URL.Query().Get("page")
	search := r.URL.Query().Get("search")

	uID, err := strconv.Atoi(id)

	if err != nil {
		fmt.Printf("could not parse user Id to get followings %s", err)
		w.WriteHeader(400)
		return
	}

	_page, err := strconv.Atoi(page)

	if err != nil {
		fmt.Printf("could not parse page to get followings %s", err)
		w.WriteHeader(400)
		return
	}

	res, err := h.store.GetUserFollowings(uID, _page, search)

	if err != nil {
		fmt.Printf("could not get user followings %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *userHandler) HandleGetUserDialogs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context())

	res, err := h.store.GetDialogs(userID)

	if err != nil {
		fmt.Printf("could not get user dialogs %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *userHandler) HandleGetDialogMessages(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		SenderID int `json:"senderId"`
		ReceiverID int `json:"receiverId"`
	} {
		SenderID: 0,
		ReceiverID: 0,
	}

	json.NewDecoder(r.Body).Decode(req)


	ms, err := h.store.GetDialogMessages(req.SenderID, req.ReceiverID)

	if err != nil {
		fmt.Printf("could not get messages of the dialog %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(ms)

}

func (h *userHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		ReceiverID int `json:"receiverId"`
		Message string `json:"message"`
	}{
		ReceiverID: 0,
		Message: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())

	err := h.store.InsertMessage(userID, req.ReceiverID, req.Message)

	if err != nil {
		fmt.Printf("could not send a message %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *userHandler) HandleAddFollower(w http.ResponseWriter, r *http.Request) {
	followerID := middleware.UserFromContext(r.Context())

	req := r.PathValue("id")

	followeeID, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("bad followee id was provided %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.AddFollower(followerID, followeeID)

	if err != nil {
		fmt.Printf("could not follow the person %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *userHandler) HandleDeleteFollow(w http.ResponseWriter, r *http.Request) {
	followerID := middleware.UserFromContext(r.Context())

	req := r.PathValue("id")

	followeeID, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("bad followee id was provided to delete follow %s", err)
		w.WriteHeader(400)
		return
	}


	err = h.store.DeleteFollow(followerID, followeeID)

	if err != nil {
		fmt.Printf("could not delete follow %s", err)
		w.WriteHeader(500)
		return
	}

}

func (h *userHandler) HandleAddFriend(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context())

	req := r.PathValue("id")

	userSId, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("bad friend id was provided %s", err)
		w.WriteHeader(400)
		return
	}

	

	err = h.store.AddFriend(userID, userSId)

	if err != nil {
		fmt.Printf("could not add friend %s", err)
		w.WriteHeader(400)
		return
	}
}

func (h *userHandler) HandleDeleteFriendship(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context())

	req := r.PathValue("id")

	userSId, err := strconv.Atoi(req)

	if err != nil {
		fmt.Printf("bad friend id was provided to delete friendship %s", err)
		w.WriteHeader(400)
		return
	}

	err = h.store.DeleteFriendship(userID, userSId)

	if err != nil {
		fmt.Printf("could not delete friendship %s", err)
		w.WriteHeader(400)
		return
	}
}

type JSONKey struct {
	Key string
	Value string
}

func (h *userHandler) HandleGetSettingsData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context())

	res, err := h.store.GetSettingsData(userID)

	if err != nil {
		fmt.Printf("could not get user settings data %s", err)
		w.WriteHeader(400)
		return
	}

	json.NewEncoder(w).Encode(res)



}

func (h *userHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		ProfileImg string `json:"profileImg"`
		BannerImg string `json:"bannerImg"`
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
		Status string `json:"status"`
	} {
		ProfileImg: "",
		BannerImg: "",
		FirstName: "",
		LastName: "",
		Status: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	userID := middleware.UserFromContext(r.Context())
	

	// acc := &struct {
	// 	UserID int
	// 	ProfileImg string 
	// } {
	// 	UserID: userID,
	// 	ProfileImg: "333",
	// }

	// data_sess, err  := json.Marshal(acc)

	// dataSToken, err := utils.CreateJWT(data_sess)

	// if err != nil {
	// 	fmt.Printf("could not create data-session cookie value %s", err)
	// 	w.WriteHeader(500)
	// 	return
	// }

	// expires := time.Now().AddDate(1, 0, 0)

	// dataSession := &http.Cookie{
	// 	Name: "data-session",
	// 	Value: dataSToken,
	// 	Expires: expires,
	// 	Path:     "/",
	// 	Secure: true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// http.SetCookie(w, dataSession)

	bytes, _ := json.Marshal(req)

	c := make(map[string]json.RawMessage)

	e := json.Unmarshal(bytes, &c)

	if e != nil {
		panic(e)
	}

	k := []JSONKey{}


	query := "update users "

	for s, j := range c {
		// fmt.Println("s", s)
		// fmt.Println("j", string(j))

		if len(j) > 2 {
			if s == "profileImg" {
				imgID := "profile-image-of-" + strconv.Itoa(userID) + "'s" + strconv.Itoa(userID) + "user"

				res, err := h.cld.Upload.Upload(context.Background(), req.ProfileImg, uploader.UploadParams{PublicID: imgID, Folder: "social-media/profile_images"});

				

				if err != nil {
					fmt.Printf("could not upload img to the cloudinary %s", err)
					return
				}

				q := JSONKey{
					Key: s,
					Value: "'" + res.SecureURL + "'",
				}

				k = append(k, q)

				acc := &struct {
					UserID int
					ProfileImg string 
				} {
					UserID: userID,
					ProfileImg: res.SecureURL,
				}

				data_sess, err  := json.Marshal(acc)

				dataSToken, err := utils.CreateJWT(data_sess)

				if err != nil {
					fmt.Printf("could not create data-session cookie value %s", err)
					w.WriteHeader(500)
					return
				}

				expires := time.Now().AddDate(1, 0, 0)

				dataSession := &http.Cookie{
					Name: "data-session",
					Value: dataSToken,
					Expires: expires,
					Path:     "/",
					Secure: true,
					SameSite: http.SameSiteNoneMode,
				}

				http.SetCookie(w, dataSession)

				continue
			}

			if s == "bannerImg" {
				imgID := "banner-image-of-" + strconv.Itoa(userID) + "'s" + strconv.Itoa(userID) + "user"

				res, err := h.cld.Upload.Upload(context.Background(), req.BannerImg, uploader.UploadParams{PublicID: imgID, Folder: "social-media/banner_images"});

				

				if err != nil {
					fmt.Printf("could not upload img to the cloudinary %s", err)
					return
				}

				q := JSONKey{
					Key: s,
					Value: "'" + res.SecureURL + "'",
				}

				k = append(k, q)

				continue
			}

			q := JSONKey{
				Key: s,
				Value: string(j),
			}

			

			k = append(k, q)
		}
		
	
	}

	for i, v := range k {
		if i == 0 {
			query += "set " + utils.CamelCaseToSnakeCase(v.Key) + " = " + utils.ReplaceDoubleQuotesWithSingle(v.Value) + ", "
		} else if i == len(k) - 1 {
			query += "" + utils.CamelCaseToSnakeCase(v.Key) + " = " + utils.ReplaceDoubleQuotesWithSingle(v.Value) + ""
		} else {
			query += "" + utils.CamelCaseToSnakeCase(v.Key) + " = " + utils.ReplaceDoubleQuotesWithSingle(v.Value) + ", "
		}

	

		
	}

	query += " where id = " + strconv.Itoa(userID)

	fmt.Println("query", query)

	err := h.store.UpdateUser(query)

	if err != nil {
		fmt.Printf("could not update the user %s", err)
		return
	}

}

func (h *userHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserFromContext(r.Context())


	fmt.Printf("user with %d connected", userID)

	conn, err := Upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: h.hub, conn: conn, send: make(chan []byte, 256), userID: userID}
	h.hub.register <- client


	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}