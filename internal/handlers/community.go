package handlers

import (
	"guthub.com/Toront0/lux-server/internal/services"

	"fmt"
	"net/http"
	"encoding/json"


)

type communityHandler struct {
	store services.CommunityStorer
}

func NewCommunityHandler(store services.CommunityStorer) *communityHandler {

	return &communityHandler{
		store: store,
	}
}

func (h *communityHandler) HandleGetCommunities(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")


	res, err := h.store.GetCommunities(search)

	if err != nil {
		fmt.Printf("could not get communities %s", err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(res)
}

