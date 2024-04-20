package handlers

import (
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"



	"net/http"
	"encoding/json"
	"fmt"
	
	"time"
)

type authHandler struct {
	store services.AuthStorer
}

func NewAuthHandler(store services.AuthStorer) *authHandler {
	return &authHandler{
		store: store,
	}
}

func (h *authHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
		Email string `json:"email"`
		Password string `json:"password"`
	} {
		FirstName: "",
		LastName: "",
		Email: "",
		Password: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	epw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("could not hash password %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acc, err := h.store.CreateUser(req.FirstName, req.LastName, req.Email, string(epw))

	if err != nil {
		fmt.Printf("could not create account %s", err)
		w.WriteHeader(400)
		return
	}

	token, err := utils.CreateJWT(acc.ID)

	if err != nil {
		fmt.Printf("could not generate JWT Token %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expires := time.Now().AddDate(1, 0, 0)

	cookie := &http.Cookie{
		Name: "jwt",
		Value: token,
		HttpOnly: true,
		Expires: expires,
		Secure: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)

	data_sess, err  := json.Marshal(acc)

	dataSToken, err := utils.CreateJWT(data_sess)

	if err != nil {
		fmt.Printf("could not create data-session cookie value %s", err)
		w.WriteHeader(500)
		return
	}

	dataSession := &http.Cookie{
		Name: "data-session",
		Value: dataSToken,
		Expires: expires,
		Secure: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, dataSession)

	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleLoginAccount(w http.ResponseWriter, r *http.Request) {
	req := &struct{
		Email string `json:"email"`
		Password string `json:"password"`
	} {
		Email: "",
		Password: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	acc, err := h.store.GetUserBy("email", req.Email)

	if err != nil {
		fmt.Printf("could not get account %s", err)
		w.WriteHeader(404)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password))

	if err != nil {
		fmt.Printf("invalid password %s", err)
		w.WriteHeader(409)
		return
	}

	token, err := utils.CreateJWT(acc.ID)

	if err != nil {
		fmt.Printf("could not generate JWT Token %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expires := time.Now().AddDate(1, 0, 0)


	// cookie := &http.Cookie{
	// 	Name: "jwt",
	// 	Value: token,
	// 	HttpOnly: true,
	// 	Expires: expires,
	// 	Path: "/",
	// 	Secure: true,
	// 	SameSite: http.SameSiteNoneMode,
	// }

	// http.SetCookie(w, cookie)


	w.Header().Add("Set-Cookie", "data-session="+dataSToken +";path=/;" + "expires=" +expires.String() + ";secure; sameSite=None; partitioned;")

	w.Header().Add("Set-Cookie", "jwt="+token +";path=/;" + "expires=" +expires.String() + ";secure;sameSite=None;partitioned;")


	data_sess, err  := json.Marshal(acc)

	dataSToken, err := utils.CreateJWT(data_sess)

	if err != nil {
		fmt.Printf("could not create data-session cookie value %s", err)
		w.WriteHeader(500)
		return
	}

	

	
	


	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	
	

	if err != nil {
		
		return
	}

	token, err := utils.ValidateJWT(cookie.Value)

	if err != nil {
		fmt.Printf("invalid JWT Token %s", err)
		w.WriteHeader(400)
		return
	}


	claims := token.Claims.(jwt.MapClaims)

	acc, _ := h.store.GetUserBy("id", claims["userID"])

	json.NewEncoder(w).Encode(acc)
}

func (h *authHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {

	c := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Secure: true,
		Expires: time.Unix(0, 0),
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}
	
	http.SetCookie(w, c)

	c = &http.Cookie{
		Name:     "data-session",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		Secure: true,
		SameSite: http.SameSiteNoneMode,
	}
	
	http.SetCookie(w, c)

}

func (h *authHandler) HandleCheckEmailExistance(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Email string `json:"email"`
	} {
		Email: "",
	}

	json.NewDecoder(r.Body).Decode(req)

	_, err := h.store.GetUserBy("email", req.Email)

	if err != nil {
		fmt.Printf("not results %s", err)
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(400)
}
