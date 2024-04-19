package api

import (
	"guthub.com/Toront0/lux-server/internal/handlers"
	"guthub.com/Toront0/lux-server/internal/handlers/user"
	"guthub.com/Toront0/lux-server/internal/services"
	"guthub.com/Toront0/lux-server/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rs/cors"
	"github.com/go-chi/chi/v5"
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

	mux := chi.NewRouter()

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

	mux.Use(c.Handler)
	mux.Use(middleware.RequireAuth)



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

	mux.Post("/sign-up", authHandler.HandleCreateAccount)
	mux.Post("/login", authHandler.HandleLoginAccount)
	mux.Get("/auth", authHandler.HandleAuthenticate)
	mux.Get("/logout", authHandler.HandleLogout)
	mux.Post("/check-email", authHandler.HandleCheckEmailExistance)

	postStore := services.NewPostStore(conn)
	postHandler := handlers.NewPostHandler(postStore, cld)

	mux.Get("/posts", postHandler.HandleGetPosts)
	mux.Post("/create-post", postHandler.HandleCreatePost)

	mux.Get("/posts/{id}", postHandler.HandleGetPost)
	mux.Get("/post-comments/{id}", postHandler.HandleGetComments)
	mux.Post("/post-comment", postHandler.HandleInsertComment)

	mux.Get("/post-comment-replies/{id}", postHandler.HandleGetCommentReplies)
	mux.Post("/post-comment-reply", postHandler.HandleInsertCommentReply)

	mux.Post("/like-post", postHandler.HandleLikePost)
	mux.Post("/delete-like-post", postHandler.HandleDeleteLike)

	mux.Post("/lpc", postHandler.HandleLikeComment)
	mux.Post("/delete-lpc", postHandler.HandleDeleteLikeComment)

	mux.Post("/lpcr", postHandler.HandleLikeCommentReply)
	mux.Post("/delete-lpcr", postHandler.HandleDeleteCommentLikeReply)

	userStore := services.NewUserStore(conn)
	userHandler := user.NewUserHandler(userStore, cld)

	mux.Get("/users/{id}", userHandler.HandleGetUserDetail)
	mux.Get("/users/{id}/posts", userHandler.HandleGetUserPosts)
	mux.Get("/users/{id}/videos", userHandler.HandleGetUserVideos)
	mux.Get("/users/{id}/friends", userHandler.HandleGetUserFriends)
	mux.Get("/users/{id}/music", userHandler.HandleGetUserMusic)
	mux.Get("/users/{id}/music-playlists", userHandler.HandleGetUserMusicPlaylists)

	mux.Post("/users/{id}/follow", userHandler.HandleAddFollower)
	mux.Post("/users/{id}/delete-follow", userHandler.HandleDeleteFollow)
	mux.Post("/users/{id}/friend", userHandler.HandleAddFriend)
	mux.Post("/users/{id}/delete-friend", userHandler.HandleDeleteFriendship)

	mux.Get("/users/{id}/followers", userHandler.HandleGetUserFollowers)
	mux.Get("/users/{id}/followings", userHandler.HandleGetUserFollowings)

	mux.Post("/users/update", userHandler.HandleUpdateUser)
	mux.Get("/settings", userHandler.HandleGetSettingsData)
	
	mux.Post("/message", userHandler.HandleSendMessage)
	mux.Post("/messages", userHandler.HandleGetDialogMessages)
	mux.Get("/dialogs", userHandler.HandleGetUserDialogs)

	mux.Get("/ws-listener", userHandler.ServeWs)

	musicStore := services.NewMusicStore(conn)
	musicHandler := handlers.NewMusicHandler(musicStore, cld)

	mux.Get("/music", musicHandler.HandleGetSongs)
	mux.Get("/playlists", musicHandler.HandleGetPlaylists)
	mux.Get("/playlists/{id}", musicHandler.HandleGetPlaylistDetail)
	mux.Post("/create-playlist", musicHandler.HandleCreatePlaylist)
	mux.Post("/delete-playlist/{id}", musicHandler.HandleDeletePlaylist)

	mux.Get("/{userId}/music/playlists/{playlistId}", musicHandler.HandleGetPlaylistSongs)
	mux.Get("/{userId}/available-songs/playlists/{playlistId}", musicHandler.HandleGetAvailableAndPlaylistSongs)
	
	mux.Post("/add-song/{id}", musicHandler.HandleAddSongToUser)
	mux.Post("/delete-song/{id}", musicHandler.HandleDeleteUserSong)
	

	videoStore := services.NewVideoStore(conn)
	videoHandler := handlers.NewVideoHandler(videoStore)

	mux.Get("/videos", videoHandler.HandleGetVideos)
	mux.Get("/videos/{id}", videoHandler.HandleGetVideo)
	mux.Post("/like-video", videoHandler.HandleLikeVideo)
	mux.Post("/delete-like-video", videoHandler.HandleDeleteLike)

	mux.Get("/video-comments/{id}", videoHandler.HandleGetComments)
	mux.Post("/video-comment", videoHandler.HandleInsertComment)
	mux.Post("/lvc", videoHandler.HandleLikeComment)
	mux.Post("/delete-lvc", videoHandler.HandleDeleteLikeComment)

	mux.Get("/video-comment-replies/{id}", videoHandler.HandleGetCommentReplies)
	mux.Post("/video-comment-reply", videoHandler.HandleInsertCommentReply)
	mux.Post("/lvcr", videoHandler.HandleLikeCommentReply)
	mux.Post("/delete-lvcr", videoHandler.HandleDeleteCommentLikeReply)

	// communityStore := services.NewCommunityStore(conn)
	// communityHandler := handlers.NewCommunityHandler(communityStore)

	// router.GET("/communities", communityHandler.HandleGetCommunities)

	// c := cors.New(cors.Options{
	// 	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins:   []string{"https://lux-client.vercel.app", "http://localhost:5173"},
	// 	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// })

	// mux := c.Handler(router)

	// stack := middleware.CreateStack(
	// 	middleware.RequireAuth,
	// )
	

	http.ListenAndServe(s.listenAddr, mux)


}