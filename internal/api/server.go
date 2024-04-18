package api

import (
	"guthub.com/Toront0/lux-server/internal/handlers"
	"guthub.com/Toront0/lux-server/internal/handlers/user"
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rs/cors"

	"net/http"

	"time"
	"context"
	"log"

)

type Post struct {
	ID int `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	AuthorName string `json:"authorName"`
	AuthorImg string `json:"authorImg"`
	Content string `json:"content"`
}

type server struct {
	listenAddr string
}

func NewServer(listenAddr string) *server {
	return &server{
		listenAddr: listenAddr,
	}
}

func (s *server) Run() {

	router := http.NewServeMux()



	conn, err := pgxpool.New(context.Background(), "postgresql://social-media_owner:P1TMr6WdiFaQ@ep-crimson-rice-a2yhkx7a.eu-central-1.aws.neon.tech/social-media?sslmode=require")
	
	
	if err != nil {
		log.Fatal("could not establish connection with DB", err)
		return
	}

	cld, err := cloudinary.NewFromParams("dwqfkmgfh", "715518655252261", "1gJkhxMKK2Rt9mSdUDF4uBqsuvg")

	if err != nil {
		log.Fatal("could not establish connection with Cloudinary CDN", err)
		return
	}

	authStore := services.NewAuthStore(conn)
	authHandler := handlers.NewAuthHandler(authStore)

	router.HandleFunc("POST /sign-up", authHandler.HandleCreateAccount)
	router.HandleFunc("POST /login", authHandler.HandleLoginAccount)
	router.HandleFunc("GET /auth", authHandler.HandleAuthenticate)
	router.HandleFunc("GET /logout", authHandler.HandleLogout)

	postStore := services.NewPostStore(conn)
	postHandler := handlers.NewPostHandler(postStore, cld)

	router.HandleFunc("GET /posts", postHandler.HandleGetPosts)
	router.HandleFunc("POST /create-post", postHandler.HandleCreatePost)

	router.HandleFunc("GET /posts/{id}", postHandler.HandleGetPost)
	router.HandleFunc("GET /post-comments/{id}", postHandler.HandleGetComments)
	router.HandleFunc("POST /post-comment", postHandler.HandleInsertComment)

	router.HandleFunc("GET /post-comment-replies/{id}", postHandler.HandleGetCommentReplies)
	router.HandleFunc("POST /post-comment-reply", postHandler.HandleInsertCommentReply)

	router.HandleFunc("POST /like-post", postHandler.HandleLikePost)
	router.HandleFunc("POST /delete-like-post", postHandler.HandleDeleteLike)

	router.HandleFunc("POST /lpc", postHandler.HandleLikeComment)
	router.HandleFunc("POST /delete-lpc", postHandler.HandleDeleteLikeComment)

	router.HandleFunc("POST /lpcr", postHandler.HandleLikeCommentReply)
	router.HandleFunc("POST /delete-lpcr", postHandler.HandleDeleteCommentLikeReply)

	userStore := services.NewUserStore(conn)
	userHandler := user.NewUserHandler(userStore, cld)

	router.HandleFunc("GET /users/{id}", userHandler.HandleGetUserDetail)
	router.HandleFunc("GET /users/{id}/posts", userHandler.HandleGetUserPosts)
	router.HandleFunc("GET /users/{id}/videos", userHandler.HandleGetUserVideos)
	router.HandleFunc("GET /users/{id}/friends", userHandler.HandleGetUserFriends)
	router.HandleFunc("GET /users/{id}/music", userHandler.HandleGetUserMusic)
	router.HandleFunc("GET /users/{id}/music-playlists", userHandler.HandleGetUserMusicPlaylists)

	router.HandleFunc("POST /users/{id}/follow", userHandler.HandleAddFollower)
	router.HandleFunc("POST /users/{id}/delete-follow", userHandler.HandleDeleteFollow)
	router.HandleFunc("POST /users/{id}/friend", userHandler.HandleAddFriend)
	router.HandleFunc("POST /users/{id}/delete-friend", userHandler.HandleDeleteFriendship)

	router.HandleFunc("GET /users/{id}/followers", userHandler.HandleGetUserFollowers)
	router.HandleFunc("GET /users/{id}/followings", userHandler.HandleGetUserFollowings)

	router.HandleFunc("POST /users/update", userHandler.HandleUpdateUser)
	router.HandleFunc("GET /settings", userHandler.HandleGetSettingsData)
	
	router.HandleFunc("POST /message", userHandler.HandleSendMessage)
	router.HandleFunc("POST /messages", userHandler.HandleGetDialogMessages)
	router.HandleFunc("GET /dialogs", userHandler.HandleGetUserDialogs)

	router.HandleFunc("GET /ws-listener/{id}", userHandler.ServeWs)

	musicStore := services.NewMusicStore(conn)
	musicHandler := handlers.NewMusicHandler(musicStore, cld)

	router.HandleFunc("GET /music", musicHandler.HandleGetSongs)
	router.HandleFunc("GET /playlists", musicHandler.HandleGetPlaylists)
	router.HandleFunc("GET /playlists/{id}", musicHandler.HandleGetPlaylistDetail)
	router.HandleFunc("POST /create-playlist", musicHandler.HandleCreatePlaylist)
	router.HandleFunc("POST /delete-playlist/{id}", musicHandler.HandleDeletePlaylist)

	router.HandleFunc("GET /{userId}/music/playlists/{playlistId}", musicHandler.HandleGetPlaylistSongs)
	router.HandleFunc("GET /{userId}/available-songs/playlists/{playlistId}", musicHandler.HandleGetAvailableAndPlaylistSongs)
	
	router.HandleFunc("POST /add-song/{id}", musicHandler.HandleAddSongToUser)
	router.HandleFunc("POST /delete-song/{id}", musicHandler.HandleDeleteUserSong)
	

	videoStore := services.NewVideoStore(conn)
	videoHandler := handlers.NewVideoHandler(videoStore)

	router.HandleFunc("GET /videos", videoHandler.HandleGetVideos)
	router.HandleFunc("GET /videos/{id}", videoHandler.HandleGetVideo)
	router.HandleFunc("POST /like-video", videoHandler.HandleLikeVideo)
	router.HandleFunc("POST /delete-like-video", videoHandler.HandleDeleteLike)

	router.HandleFunc("GET /video-comments/{id}", videoHandler.HandleGetComments)
	router.HandleFunc("POST /video-comment", videoHandler.HandleInsertComment)
	router.HandleFunc("POST /lvc", videoHandler.HandleLikeComment)
	router.HandleFunc("POST /delete-lvc", videoHandler.HandleDeleteLikeComment)

	router.HandleFunc("GET /video-comment-replies/{id}", videoHandler.HandleGetCommentReplies)
	router.HandleFunc("POST /video-comment-reply", videoHandler.HandleInsertCommentReply)
	router.HandleFunc("POST /lvcr", videoHandler.HandleLikeCommentReply)
	router.HandleFunc("POST /delete-lvcr", videoHandler.HandleDeleteCommentLikeReply)

	// communityStore := services.NewCommunityStore(conn)
	// communityHandler := handlers.NewCommunityHandler(communityStore)

	// router.HandleFunc("GET /communities", communityHandler.HandleGetCommunities)




	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	mux := c.Handler(router)

	stack := middleware.CreateStack(
		middleware.RequireAuth,
	)
	

	server := http.Server{
		Addr: s.listenAddr,
		Handler: stack(mux),
	}

	server.ListenAndServe()


}